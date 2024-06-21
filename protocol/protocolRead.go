package protocol

import (
	"errors"
	"main/config"
	"main/utils"
	"strconv"
	"strings"
)

func ReadSimpleString(input string) ([]string, string, error) {
	endIndex := strings.Index(input, config.CLRF)

	if endIndex == -1 {
		return []string{}, "", errors.New("Can read simple string wrong data: " + input)
	}

	return []string{input[1:endIndex]}, input[endIndex+len(config.CLRF):], nil
}

func ReadBulkString(input string) ([]string, string, error) {
	lengthEndIndex := strings.Index(input, config.CLRF)

	if lengthEndIndex == -1 {
		return []string{}, "", errors.New("Can read simple string wrong data: " + input)
	}

	length, err := strconv.Atoi(input[1:lengthEndIndex])

	if err != nil {
		return []string{}, "", errors.New("Couldnt convert input length: " + input)
	}

	wordStartIndex := lengthEndIndex + len(config.CLRF)
	wordEndIndex := wordStartIndex + length

	return []string{input[wordStartIndex:wordEndIndex]}, input[wordEndIndex+len(config.CLRF):], nil
}

// *2\r\n$5\r\nhello\r\n$5\r\nworld\r\n
func ReadArray(input string) ([]string, string, error) {
	if len(input) < 2 {
		return []string{}, "", errors.New("Data need to have at least lenth of 2: " + input)
	}

	endOfLength := strings.Index(input, config.CLRF)
	if endOfLength == -1 {
		return []string{}, "", errors.New("Couldnt find CLRF for array length on input: " + input)
	}

	arrayLength, err := strconv.Atoi(input[1:endOfLength])
	if err != nil {
		return []string{}, "", errors.New("Couldnt convert number of array element to number on input: " + input)
	}

	numberOfCLRF := arrayLength*2 + 1
	arrayInputEndIndex := utils.FindNOccurance(input, config.CLRF, numberOfCLRF)
	arrayInput := input[:arrayInputEndIndex+len(config.CLRF)]

	items := strings.Split(arrayInput, config.CLRF)

	resp := make([]string, 0, arrayLength)

	for i := 2; i < len(items); i += 2 {
		resp = append(resp, items[i])
	}

	return resp, input[arrayInputEndIndex+len(config.CLRF):], nil
}

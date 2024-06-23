package protocol

import (
	"errors"
	"fmt"
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

	return []string{input[1:endIndex]}, input[endIndex+config.CLRFLength:], nil
}

// $5\r\nhello\r\n
func ReadBulkString(input string) ([]string, string, error) {
	msgLengthEndIndex := 1
	for input[msgLengthEndIndex] != byte(config.CLRF[0]) {
		if input[msgLengthEndIndex] < '0' || input[msgLengthEndIndex] > '9' {
			return []string{}, "", &utils.AppError{
				ErrType: utils.InvalidCharInMsgLength,
				Msg:     fmt.Sprintf("Char: %q, its not a number", input[msgLengthEndIndex]),
			}
		}
		if msgLengthEndIndex < len(input)-1 {
			msgLengthEndIndex += 1
		} else {
			return []string{}, input, nil
		}
	}

	stringNumber := input[1:msgLengthEndIndex]
	msgLength, err := strconv.Atoi(stringNumber)

	if err != nil {
		return []string{}, "", &utils.AppError{
			ErrType: utils.CannotConvertStringToNumber,
			Msg:     fmt.Sprintf("Cant conver't:%q to number", stringNumber),
		}
	}

	wordStartIndex := msgLengthEndIndex + len(config.CLRF)
	wordEndIndex := wordStartIndex + msgLength
	if len(input) < wordEndIndex {
		return []string{}, input, nil
	}

	word := input[wordStartIndex:wordEndIndex]
	if len(input) < wordEndIndex+config.CLRFLength {
		return []string{}, input, nil
	}

	endCommandIndex := wordEndIndex + config.CLRFLength

	if input[wordEndIndex:endCommandIndex] != config.CLRF {
		return []string{}, "", &utils.AppError{
			ErrType: utils.WrongCommandFormat,
			Msg:     fmt.Sprintf("Command should end up with: %q, insted we got: %q", config.CLRF, input[wordEndIndex:endCommandIndex]),
		}
	}
	fmt.Println("hello we read")
	return []string{word}, input[endCommandIndex:], nil
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
	arrayInput := input[:arrayInputEndIndex+config.CLRFLength]

	items := strings.Split(arrayInput, config.CLRF)

	resp := make([]string, 0, arrayLength)

	for i := 2; i < len(items); i += 2 {
		resp = append(resp, items[i])
	}

	return resp, input[arrayInputEndIndex+config.CLRFLength:], nil
}

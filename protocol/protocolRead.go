package protocol

import (
	"fmt"
	"main/config"
	"main/utils"
	"strconv"
)

func ReadBulkString(input string) ([]string, string, error) {
	if len(input) <= 1 {
		return []string{}, input, nil
	}

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

	if wordStartIndex > len(input) {
		return []string{}, input, nil
	} else if input[msgLengthEndIndex:wordStartIndex] != config.CLRF {
		return []string{}, "", &utils.AppError{
			ErrType: utils.WrongCommandFormat,
			Msg:     fmt.Sprintf("Wrong format, expected CLRF after msg length for input: %q", input),
		}
	}

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

	return []string{word}, input[endCommandIndex:], nil
}

func ReadArray(input string) ([]string, string, error) {

	if len(input) <= 1 {
		return []string{}, input, nil
	}

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
	arrayLength, err := strconv.Atoi(stringNumber)

	if err != nil {
		return []string{}, "", &utils.AppError{
			ErrType: utils.CannotConvertStringToNumber,
			Msg:     fmt.Sprintf("Cant convert: %q to number", stringNumber),
		}
	}

	bulkStringStartIndex := msgLengthEndIndex + config.CLRFLength

	if bulkStringStartIndex > len(input) {
		return []string{}, input, nil
	} else if input[msgLengthEndIndex:bulkStringStartIndex] != config.CLRF {
		return []string{}, "", &utils.AppError{
			ErrType: utils.WrongCommandFormat,
			Msg:     fmt.Sprintf("Wrong format, expected CLRF after msg length for input: %q", input),
		}
	}

	currentInput := input[bulkStringStartIndex:]
	resp := []string{}

	for i := 0; i < arrayLength; i++ {
		result, rest, err := ReadBulkString(currentInput)

		if err != nil {
			return []string{}, "", err
		}

		resp = append(resp, result...)
		currentInput = rest
	}

	return resp, currentInput, nil
}

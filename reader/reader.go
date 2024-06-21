package reader

import (
	"errors"
	"fmt"
	"io"
	"main/protocol"
)

func ProcessData(reader io.Reader) ([]string, error) {
	data, err := Read(reader)

	if err != nil {
		return []string{}, err
	}

	processedData, err := process(data)

	if err != nil {
		return []string{}, err
	}

	return processedData, nil
}

func Read(reader io.Reader) (string, error) {
	buff := make([]byte, 1024)
	n, err := reader.Read(buff)

	if err != nil {
		fmt.Errorf("Problem reading data from connection, err:%q", err)
		return "", err
	}

	data := buff[:n]

	return string(data), nil
}

func process(data string) (res []string, err error) {

	if len(data) == 0 {
		fmt.Print("no data to process")
		return []string{}, nil
	}

	currentInput := data

	for currentInput != "" {
		operator := currentInput[0]
		switch operator {
		case '$':
			res, currentInput, err = protocol.ReadBulkString(currentInput)
		case '*':
			res, currentInput, err = protocol.ReadArray(currentInput)
		case '+':
			res, currentInput, err = protocol.ReadSimpleString(currentInput)
		default:
			fmt.Errorf("Operator: %v not supported", operator)
			err = errors.New("Operation not supported")
		}

		if err != nil {
			break
		}
	}

	return res, err
}

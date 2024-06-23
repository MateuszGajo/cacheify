package reader

import (
	"errors"
	"fmt"
	"io"
	"main/protocol"
)

func ProcessData(reader io.Reader) ([][]string, error) {
	data, err := Read(reader)

	if err != nil {
		return [][]string{}, err
	}

	processedData, err := process(data)

	if err != nil {
		return [][]string{}, err
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

func process(data string) (resp [][]string, err error) {

	if len(data) == 0 {
		fmt.Print("no data to process")
		return [][]string{}, nil
	}

	currentInput := data

	for currentInput != "" {
		operator := currentInput[0]
		commands := []string{}
		switch operator {
		case '$':
			commands, currentInput, err = protocol.ReadBulkString(currentInput)
		case '*':
			commands, currentInput, err = protocol.ReadArray(currentInput)
		default:
			fmt.Errorf("Operator: %v not supported", operator)
			err = errors.New("Operation not supported")
		}
		resp = append(resp, commands)

		if err != nil {
			return resp, err
		}
	}

	return resp, err
}

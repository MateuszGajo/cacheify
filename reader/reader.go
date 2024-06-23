package reader

import (
	"errors"
	"fmt"
	"io"
	"main/protocol"
)

func ProcessData(reader io.Reader) ([][]string, string, error) {
	data, err := Read(reader)

	if err != nil {
		return [][]string{}, "", err
	}

	processedData, rest, err := process(data)

	if err != nil {
		return [][]string{}, "", err
	}

	return processedData, rest, nil
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

func process(data string) (resp [][]string, unprocessData string, err error) {

	if len(data) == 0 {
		fmt.Print("no data to process")
		return [][]string{}, unprocessData, nil
	}

	unprocessData = data
	cuttedData := ""

	for unprocessData != "" {
		operator := unprocessData[0]
		commands := []string{}
		switch operator {
		case '$':
			commands, cuttedData, unprocessData, err = protocol.ReadBulkString(unprocessData)
		case '*':
			commands, cuttedData, unprocessData, err = protocol.ReadArray(unprocessData)
		default:
			fmt.Errorf("Operator: %v not supported", operator)
			err = errors.New("Operation not supported")
		}
		resp = append(resp, commands)

		if err != nil {
			return resp, cuttedData, err
		}
		if cuttedData != "" {
			return resp, cuttedData, err
		}
	}

	return resp, cuttedData, err
}

package reader

import (
	"errors"
	"fmt"
	"io"
	"main/protocol"
)

var dataUnprocess string

func ProcessData(data string) ([][]string, error) {

	input := dataUnprocess + data

	processedData, rest, err := process(input)

	if err != nil {
		return [][]string{}, err
	}
	dataUnprocess = rest

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

func findFirstOccurance(input string, delimiters []byte) int {
	for i := 0; i < len(input); i++ {
		for j := 0; j < len(delimiters); j++ {
			if input[i] == delimiters[j] {
				return i
			}
		}
	}
	return -1
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
			// If error find start of next command or if there is not, discard all corurpted data
			index := findFirstOccurance(unprocessData, []byte{36, 42})
			if index == -1 {
				unprocessData = ""
			} else {
				unprocessData = unprocessData[index:]
			}
		}
		if cuttedData != "" {
			return resp, cuttedData, err
		}
	}

	return resp, cuttedData, err
}

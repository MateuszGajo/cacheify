package reader

import (
	"errors"
	"fmt"
	"io"
	"main/protocol"
)

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

func Process(data string) (res []string, err error) {

	if len(data) == 0 {
		fmt.Print("no data to process")
		return []string{}, nil
	}

	operator := data[0]

	switch operator {
	case '$':
		res, data, err = protocol.ReadSimpleString(data)

	default:
		fmt.Errorf("Operator: %v not supported", operator)
		err = errors.New("Operation not supported")
	}

	return res, err
}

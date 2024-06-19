package protocol

import (
	"errors"
	"main/config"
	"strings"
)

func ReadSimpleString(input string) ([]string, string, error) {
	endIndex := strings.Index(input, config.CLRF)

	if endIndex == -1 {
		return []string{}, "", errors.New("Can read simple string wrong data: " + input)
	}

	return []string{input[1:endIndex]}, input[endIndex+len(config.CLRF):], nil
}

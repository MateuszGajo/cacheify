package protocol

import (
	"fmt"
	"main/config"
)

func WriteSimpleString(msg string) string {
	return fmt.Sprintf("+%v%v", msg, config.CLRF)
}

func WriteBulkString(msg string) string {
	return fmt.Sprintf("$%v%v%v%v", len(msg), config.CLRF, msg, config.CLRF)
}

func WriteArrayString(msg []string) string {
	input := fmt.Sprintf("*%v%v", len(msg), config.CLRF)

	for i := 0; i < len(msg); i++ {
		input += fmt.Sprintf("$%v%v%v%v", len(msg[i]), config.CLRF, msg[i], config.CLRF)
	}

	return input
}

func WriteSimpleError(msg string) string {
	return fmt.Sprintf("-%v%v", msg, config.CLRF)
}

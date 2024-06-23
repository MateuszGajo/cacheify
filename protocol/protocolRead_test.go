package protocol

import (
	"main/utils"
	"reflect"
	"testing"
)

func TestBulkString(t *testing.T) {
	input := "$5\r\nhello\r\n"
	resp, _, _, _ := ReadBulkString(input)

	expectedRes := "hello"

	if resp[0] != expectedRes {
		t.Fatalf("Expected:%q, got:%q", expectedRes, resp)
	}
}

func TestBulkInvalidMsgLongerThanDeclared(t *testing.T) {
	input := "$5\r\nhelloa\r\n"
	resp, _, _, err := ReadBulkString(input)

	if err == nil {
		t.Fatalf("Should return error for input:%q, we got res :%q", input, resp)
	}
}

func TestBulkInvalidMsgLength(t *testing.T) {
	input := "$5a\r\nhelloa\r\n"
	_, _, _, err := ReadBulkString(input)

	expectedErrType := utils.InvalidCharInMsgLength
	errType := utils.GetErrorType(err)

	if !(errType == expectedErrType) {
		t.Fatalf("expected error type: %v, got err:%v", expectedErrType, errType)
	}
}

func TestBulkWithJustOneChar(t *testing.T) {
	input := "$"
	_, cuttedData, _, _ := ReadBulkString(input)

	expectedRes := "$"

	if !(cuttedData == expectedRes) {
		t.Fatalf("expected error type: %v, got err:%v", expectedRes, cuttedData)
	}
}

func TestBulkInvalidEndingDelimiter(t *testing.T) {
	input := "$5\r\nhello\r\\2"
	_, _, _, err := ReadBulkString(input)

	expectedErrType := utils.WrongCommandFormat
	errType := utils.GetErrorType(err)

	if !(errType == expectedErrType) {
		t.Fatalf("expected error type: %v, got err:%v", expectedErrType, errType)
	}

}

func TestBulkNotFullMessage(t *testing.T) {
	input := "$5\r\nhello"
	_, cuttedData, _, _ := ReadBulkString(input)

	expectedRes := "$5\r\nhello"

	if cuttedData != expectedRes {
		t.Fatalf("Expected:%q, got:%q", expectedRes, cuttedData)
	}
}

func TestBulkWrondDelimiterAfterMsgLength(t *testing.T) {
	input := "$5\rchello"
	_, _, _, err := ReadBulkString(input)

	expectedErrType := utils.WrongCommandFormat
	errType := utils.GetErrorType(err)

	if !(errType == expectedErrType) {
		t.Fatalf("expected error type: %v, got err:%v", expectedErrType, errType)
	}
}

func TestBulkStringWithAnotherCommand(t *testing.T) {
	input := "$5\r\nhello\r\n+OK\r\n"
	_, _, unprocessData, _ := ReadBulkString(input)

	expectedRes := "+OK\r\n"

	if unprocessData != expectedRes {
		t.Fatalf("Expected:%q, got:%q", expectedRes, unprocessData)
	}
}

func TestArray(t *testing.T) {
	input := "*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n"
	resp, _, _, _ := ReadArray(input)

	if !reflect.DeepEqual(resp, []string{"hello", "world"}) {
		t.Fatalf("Expected:%q, got:%q", []string{"hello", "world"}, resp)
	}
}

func TestArrayWithAnotherCommand(t *testing.T) {
	input := "*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n+OK\r\n+OK\r\n"
	_, _, unprocessData, _ := ReadArray(input)

	if unprocessData != "+OK\r\n+OK\r\n" {
		t.Fatalf("Expected:%q, got:%q", "+OK\r\n+OK\r\n", unprocessData)
	}

}

func TestArrayWithInvalidMsgLength(t *testing.T) {
	input := "*2g\r\n$5\r\nhello\r\n$5\r\nworld\r\n+OK\r\n+OK\r\n"
	_, _, _, err := ReadArray(input)

	expectedErrType := utils.InvalidCharInMsgLength
	errType := utils.GetErrorType(err)

	if !(errType == expectedErrType) {
		t.Fatalf("expected error type: %v, got err:%v", expectedErrType, errType)
	}

}

func TestArrayWithInvalidDelimiterAfterMsgLength(t *testing.T) {
	input := "*2\rc$5\r\nhello\r\n$5\r\nworld\r\n+OK\r\n+OK\r\n"
	_, _, _, err := ReadArray(input)

	expectedErrType := utils.WrongCommandFormat
	errType := utils.GetErrorType(err)

	if !(errType == expectedErrType) {
		t.Fatalf("expected error type: %v, got err:%v", expectedErrType, errType)
	}

}

func TestArrayWithErrorInsideBulkString(t *testing.T) {
	input := "*2\r\n$5insideBulkdString\r\nhello\r\n$5\r\nworld\r\n+OK\r\n+OK\r\n"
	_, _, _, err := ReadArray(input)

	expectedErrType := utils.InvalidCharInMsgLength
	errType := utils.GetErrorType(err)

	if !(errType == expectedErrType) {
		t.Fatalf("expected error type: %v, got err:%v", expectedErrType, errType)
	}

}

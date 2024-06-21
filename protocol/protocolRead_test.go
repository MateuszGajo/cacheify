package protocol

import (
	"reflect"
	"testing"
)

func TestSimpleString(t *testing.T) {
	input := "+OK\r\n"
	resp, newInput, _ := ReadSimpleString(input)

	if resp[0] != "OK" {
		t.Fatalf("Expected:%q, got:%q", "OK", resp)
	}

	if newInput != "" {
		t.Fatalf("Exepected new input to be empty insted we got: %q", newInput)
	}
}

func TestSimpleStringWithAnotherCommand(t *testing.T) {
	input := "+OK\r\n+OK\r\n"
	resp, newInput, _ := ReadSimpleString(input)

	if resp[0] != "OK" {
		t.Fatalf("Expected:%q, got:%q", "OK", resp)
	}

	if newInput != "+OK\r\n" {
		t.Fatalf("Exepected new input to be:%q, got: %q", "+OK\r\n", newInput)
	}
}

func TestSimpleInvalidString(t *testing.T) {
	input := "+OK\\n"
	resp, _, err := ReadSimpleString(input)

	if err == nil {
		t.Fatalf("Should return error insted we got response :%q", resp)
	}
}

func TestBulkString(t *testing.T) {
	input := "$5\r\nhello\r\n"
	resp, _, _ := ReadBulkString(input)

	expectedRes := "hello"

	if resp[0] != expectedRes {
		t.Fatalf("Expected:%q, got:%q", expectedRes, resp)
	}
}

func TestBulkStringWithAnotherCommand(t *testing.T) {
	input := "$5\r\nhello\r\n+OK\r\n"
	_, rest, _ := ReadBulkString(input)

	expectedRes := "+OK\r\n"

	if rest != expectedRes {
		t.Fatalf("Expected:%q, got:%q", expectedRes, rest)
	}
}

func TestArray(t *testing.T) {
	input := "*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n"
	resp, _, _ := ReadArray(input)

	if !reflect.DeepEqual(resp, []string{"hello", "world"}) {
		t.Fatalf("Expected:%q, got:%q", []string{"hello", "world"}, resp)
	}
}

func TestArrayWithAnotherCommand(t *testing.T) {
	input := "*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n+OK\r\n+OK\r\n"
	_, rest, _ := ReadArray(input)

	if rest != "+OK\r\n+OK\r\n" {
		t.Fatalf("Expected:%q, got:%q", "+OK\r\n+OK\r\n", rest)
	}

}

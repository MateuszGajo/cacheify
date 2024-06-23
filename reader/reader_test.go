package reader

import (
	"bytes"
	"reflect"
	"testing"
)

func TestProcessSingleCommand(t *testing.T) {
	input := []byte("*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n")
	res, _, _ := ProcessData(bytes.NewReader(input))

	expectedRes := []string{"hello", "world"}

	if !reflect.DeepEqual(expectedRes, res[0]) {
		t.Errorf("expected:%q, got:%q", expectedRes, res[0])
	}
}

func TestProcessMultipleCommand(t *testing.T) {
	input := []byte("*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n$3\r\nxyz\r\n")
	res, _, _ := ProcessData(bytes.NewReader(input))

	expectedRes := "xyz"

	if res[1][0] != expectedRes {
		t.Errorf("expected:%q, got:%q", expectedRes, res[1][0])
	}
}

func TestProcessMultipleCommandWithOneIcomplete(t *testing.T) {
	input := []byte("*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n$3\r\nxyz\r\n$3\r")
	_, rest, _ := ProcessData(bytes.NewReader(input))

	expectedRes := "$3\r"

	if rest != expectedRes {
		t.Errorf("expected:%q, got:%q", expectedRes, rest)
	}
}

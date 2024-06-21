package reader

import (
	"bytes"
	"reflect"
	"testing"
)

func TestProcessSingleCommand(t *testing.T) {
	input := []byte("*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n")
	res, _ := ProcessData(bytes.NewReader(input))

	expectedRes := []string{"hello", "world"}

	if !reflect.DeepEqual(expectedRes, res) {
		t.Errorf("expected:%q, got:%q", expectedRes, res)
	}
}

func TestProcessMultipleCommand(t *testing.T) {
	input := []byte("*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n$3\r\nxyz\r\n")
	res, err := ProcessData(bytes.NewReader(input))

	t.Error(res, err)

	expectedRes := "xyz"

	if res[2] != expectedRes {
		t.Errorf("expected:%q, got:%q", expectedRes, res)
	}
}

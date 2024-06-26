package reader

import (
	"reflect"
	"testing"
)

func TestProcessSingleCommand(t *testing.T) {
	input := "*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n"
	res, _ := ProcessData(input)

	expectedRes := []string{"hello", "world"}

	if !reflect.DeepEqual(expectedRes, res[0]) {
		t.Errorf("expected:%q, got:%q", expectedRes, res[0])
	}
}

func TestProcessMultipleCommand(t *testing.T) {
	input := "*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n$3\r\nxyz\r\n"
	res, _ := ProcessData(input)

	expectedRes := "xyz"

	if res[1][0] != expectedRes {
		t.Errorf("expected:%q, got:%q", expectedRes, res[1][0])
	}
}

func TestProcessCommandInMultipleInvocation(t *testing.T) {
	inputFirst := "$3\r\n"
	inputSecond := "xyz\r\n"
	ProcessData(inputFirst)
	res, _ := ProcessData(inputSecond)

	expectedRes := "xyz"

	if res[0][0] != expectedRes {
		t.Errorf("expected:%q, got:%q", expectedRes, res[1][0])
	}
}

// func TestProcessMultipleCommandWithOneIcomplete(t *testing.T) {
// 	input := "*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n$3\r\nxyz\r\n$3\r"
// 	_, rest := ProcessData(input)

// 	expectedRes := "$3\r"

// 	if rest != expectedRes {
// 		t.Errorf("expected:%q, got:%q", expectedRes, rest)
// 	}
// }

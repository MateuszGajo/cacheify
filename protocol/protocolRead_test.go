package protocol

import "testing"

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

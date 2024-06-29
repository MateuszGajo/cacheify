package protocol

import "testing"

func TestSimpleWrite(t *testing.T) {
	resp := WriteSimpleString("OK")

	expected := "+OK\r\n"
	if resp != expected {
		t.Fatalf("Expected: %q, got:%q", expected, resp)
	}
}

func TestBulkdWrite(t *testing.T) {
	resp := WriteBulkString("abc")

	expected := "$3\r\nabc\r\n"
	if resp != expected {
		t.Fatalf("Expected: %q, got:%q", expected, resp)
	}
}

func TestWriteArrayString(t *testing.T) {
	resp := WriteArrayString([]string{"ECHO", "hey"})

	expected := "*2\r\n$4\r\nECHO\r\n$3\r\nhey\r\n"
	if resp != expected {
		t.Fatalf("Expected: %q, got:%q", expected, resp)
	}
}

func TestWriteSimpleError(t *testing.T) {
	resp := WriteSimpleError("unknown command")

	expected := "-unknown command\r\n"
	if resp != expected {
		t.Fatalf("Expected: %q, got:%q", expected, resp)
	}
}

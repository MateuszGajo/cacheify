package protocol

import "testing"

func TestSimpleWrite(t *testing.T) {
	resp := WriteSimpleString("OK")

	expected := "+OK\r\n"
	if resp != expected {
		t.Fatalf("Expected: %v, got:%v", expected, resp)
	}
}

func TestWriteArrayString(t *testing.T) {
	resp := WriteArrayString([]string{"ECHO", "hey"})

	expected := "*2\r\n$4\r\nECHO\r\n$3\r\nhey\r\n"
	if resp != expected {
		t.Fatalf("Expected: %q, got:%q", expected, resp)
	}
}

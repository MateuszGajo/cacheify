package server

import (
	"main/protocol"
	"reflect"
	"testing"
)

func TestEchoComman(t *testing.T) {
	address := "127.0.0.1"
	port := "6380"
	server := createAndStartServer(address, port)
	conn := connectToServer("tcp", address+":"+port, t)

	conn.Write([]byte(protocol.WriteArrayString([]string{"ECHO", "key"})))

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)

	if err != nil {
		t.Fatalf("problem reading: %v", err)
	}

	msg := buf[:n]

	expected := []byte(protocol.WriteSimpleString("OK"))

	if !reflect.DeepEqual(msg, expected) {
		t.Fatalf("Expected: %q, got: %q", string(expected), string(msg))
	}

	defer func() {
		server.Close()
	}()
}

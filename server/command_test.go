package server

import (
	"main/protocol"
	"reflect"
	"testing"
	"time"
)

func TestPingCommand(t *testing.T) {
	address := "127.0.0.1"
	port := "6380"
	server := createAndStartServer(address, port)
	conn := connectToServer("tcp", address+":"+port, t)

	conn.Write([]byte(protocol.WriteArrayString([]string{"PING"})))

	msg := readData(conn, t)
	expected := protocol.WriteSimpleString("PONG")

	if !reflect.DeepEqual(msg, expected) {
		t.Fatalf("Expected: %q, got: %q", expected, msg)
	}

	defer func() {
		server.Close()
		conn.Close()
	}()
}

func TestEchoCommand(t *testing.T) {
	address := "127.0.0.1"
	port := "6380"
	server := createAndStartServer(address, port)
	conn := connectToServer("tcp", address+":"+port, t)

	conn.Write([]byte(protocol.WriteArrayString([]string{"ECHO", "key"})))

	msg := readData(conn, t)
	expected := protocol.WriteBulkString("key")

	if !reflect.DeepEqual(msg, expected) {
		t.Fatalf("Expected: %q, got: %q", expected, msg)
	}

	defer func() {
		server.Close()
		conn.Close()
	}()
}

func TestSetCommand(t *testing.T) {
	address := "127.0.0.1"
	port := "6380"
	server := createAndStartServer(address, port)
	conn := connectToServer("tcp", address+":"+port, t)

	conn.Write([]byte(protocol.WriteArrayString([]string{"SET", "abc", "def"})))

	msg := readData(conn, t)
	expected := protocol.WriteSimpleString("OK")

	if !reflect.DeepEqual(msg, expected) {
		t.Fatalf("Expected: %q, got: %q", expected, msg)
	}

	defer func() {
		server.Close()
		conn.Close()
	}()
}

func TestSetNoArgs(t *testing.T) {
	address := "127.0.0.1"
	port := "6380"
	server := createAndStartServer(address, port)
	conn := connectToServer("tcp", address+":"+port, t)

	conn.Write([]byte(protocol.WriteArrayString([]string{"SET"})))

	msg := readData(conn, t)
	expected := protocol.WriteSimpleError("We need at least key and value in set command")

	if !reflect.DeepEqual(msg, expected) {
		t.Fatalf("Expected: %q, got: %q", expected, msg)
	}

	defer func() {
		server.Close()
		conn.Close()
	}()
}

func TestGetCommandNoValue(t *testing.T) {
	address := "127.0.0.1"
	port := "6380"
	server := createAndStartServer(address, port)
	conn := connectToServer("tcp", address+":"+port, t)

	conn.Write([]byte(protocol.WriteArrayString([]string{"GET", "abc"})))

	msg := readData(conn, t)
	expected := protocol.WriteBulkString("-1")

	if !reflect.DeepEqual(msg, expected) {
		t.Fatalf("Expected: %q, got: %q", expected, msg)
	}

	defer func() {
		server.Close()
		conn.Close()
	}()
}

func TestGetCommand(t *testing.T) {
	address := "127.0.0.1"
	port := "6380"
	server := createAndStartServer(address, port)
	conn := connectToServer("tcp", address+":"+port, t)

	key := "abc"
	value := "def"
	setValue(conn, key, value, t)

	conn.Write([]byte(protocol.WriteArrayString([]string{"GET", key})))

	msg := readData(conn, t)
	expected := protocol.WriteBulkString(value)

	if !reflect.DeepEqual(msg, expected) {
		t.Fatalf("Expected: %q, got: %q", expected, msg)
	}

	defer func() {
		server.Close()
		conn.Close()
	}()
}

func TestSetCommandWithExpiry(t *testing.T) {
	address := "127.0.0.1"
	port := "6380"
	server := createAndStartServer(address, port)
	conn := connectToServer("tcp", address+":"+port, t)

	key := "abc"
	value := "def"
	expTimeMs := 100

	setValueWithExp(conn, key, value, PX, expTimeMs, t)

	msg := getData(conn, key, t)
	expected := protocol.WriteBulkString(value)

	if msg != expected {
		t.Fatalf("Value should be:%q, insted we got:%q", expected, value)
	}

	time.Sleep(time.Duration(expTimeMs) * time.Millisecond)

	msg = getData(conn, key, t)
	expected = protocol.WriteBulkString("-1")

	if msg != expected {
		t.Fatalf("Value should be:%q, insted we got:%q", expected, value)
	}

	defer func() {
		server.Close()
		conn.Close()
	}()

}

func TestInvalidSetCommandWithExpiry(t *testing.T) {
	address := "127.0.0.1"
	port := "6380"
	server := createAndStartServer(address, port)
	conn := connectToServer("tcp", address+":"+port, t)

	key := "abc"
	value := "def"

	conn.Write([]byte(protocol.WriteArrayString([]string{"SET", key, value, string(PX), "aa"})))
	msg := readData(conn, t)
	expected := protocol.WriteSimpleError("Invalid expiry time")

	if msg != expected {
		t.Fatalf("Expected:%q, got:%q", expected, msg)
	}

	conn.Write([]byte(protocol.WriteArrayString([]string{"SET", key, value, string(PX), "-1"})))
	msg = readData(conn, t)
	expected = protocol.WriteSimpleError("Expiry time need to be greater equal 0")

	if msg != expected {
		t.Fatalf("Expected:%q, got:%q", expected, msg)
	}

	defer func() {
		server.Close()
		conn.Close()
	}()

}

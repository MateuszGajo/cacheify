package server

import (
	"io"
	"main/protocol"
	"net"
	"reflect"
	"strconv"
	"sync"
	"testing"
)

func tempIsConnClosed(conn net.Conn) bool {
	buf := make([]byte, 1)

	_, err := conn.Read(buf)

	return err == io.EOF
}

type Resp struct {
	conn net.Conn
	err  error
}

func connectToServer(network, address string, t *testing.T) net.Conn {
	connChan := make(chan Resp)

	go func() {
		conn, err := net.Dial(network, address)
		connChan <- Resp{conn: conn, err: err}
	}()

	resp := <-connChan

	if resp.err != nil {
		t.Fatalf("couldnt connect to server :%v", resp.err)
	}

	return resp.conn
}

func createAndStartServer(address, port string) *Server {
	server := CreateServer(
		WithPort(port),
		WithAddress(address),
	)
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		wg.Done()
		server.RunServer()
	}()

	wg.Wait()

	return server
}

func readData(conn net.Conn, t *testing.T) string {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)

	if err != nil {
		t.Fatalf("problem reading: %v", err)
	}

	msg := buf[:n]

	return string(msg)
}

func setValue(conn net.Conn, key, value string, t *testing.T) {
	conn.Write([]byte(protocol.WriteArrayString([]string{"SET", key, value})))
	msg := readData(conn, t)

	expected := protocol.WriteSimpleString("OK")

	if !reflect.DeepEqual(msg, expected) {
		t.Fatalf("Expected: %q, got: %q", expected, msg)
	}
}

func setValueWithExp(conn net.Conn, key, value string, expType ExpType, expValueMs int, t *testing.T) {
	conn.Write([]byte(protocol.WriteArrayString([]string{"SET", key, value, string(expType), strconv.Itoa(expValueMs)})))
	msg := readData(conn, t)

	expected := protocol.WriteSimpleString("OK")

	if !reflect.DeepEqual(msg, expected) {
		t.Fatalf("Expected: %q, got: %q", expected, msg)
	}
}

func getData(conn net.Conn, key string, t *testing.T) string {

	conn.Write([]byte(protocol.WriteArrayString([]string{"GET", key})))

	msg := readData(conn, t)

	return msg
}

func TestTerminationConnectionShutdown(t *testing.T) {
	address := "127.0.0.1"
	port := "6380"
	server := createAndStartServer(address, port)
	conn := connectToServer("tcp", address+":"+port, t)

	server.Close()

	if !tempIsConnClosed(conn) {
		t.Fatal("Connection should be terminated")
	}

	defer func() {
		conn.Close()
	}()

}

package server

import (
	"io"
	"net"
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
	server := CreateServer(address, port)
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		wg.Done()
		server.RunServer()
	}()

	wg.Wait()

	return server
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

}

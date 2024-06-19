package server

import (
	"fmt"
	"net"
	"os"
	"sync"
)

type Server struct {
	listener net.Listener
	quit     chan interface{}
	wg       sync.WaitGroup
}

type Command string

const (
	Echo Command = "ECHO"
)

func CreateServer(address, port string) *Server {
	fmt.Println("what we running on")

	server := &Server{
		quit: make(chan interface{}),
	}

	listener, err := net.Listen("tcp", address+":"+port)

	if err != nil {
		fmt.Print("can't run server", err)
		os.Exit(1)
	}

	server.listener = listener
	server.wg.Add(1)

	return server
}

func (server *Server) RunServer() {

	defer server.wg.Done()

	for {
		chanConn := make(chan net.Conn, 1)
		go func() {
			conn, err := server.listener.Accept()

			if err != nil {
				fmt.Print("problem with accepting connection", err)
				return
			}
			chanConn <- conn
		}()
		select {
		case <-server.quit:
			return
		case a := <-chanConn:
			server.wg.Add(1)
			go func() {
				handleConn(a)
				server.wg.Done()
			}()
		}
	}
}

func (server *Server) Close() {
	close(server.quit)
	err := server.listener.Close()
	if err != nil {
		fmt.Print("error ocurred while closing the server")
	}
	server.wg.Wait()
}

func handleConn(conn net.Conn) {

}

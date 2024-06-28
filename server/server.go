package server

import (
	"errors"
	"fmt"
	"io"
	"main/protocol"
	"main/reader"
	"net"
	"os"
	"strings"
	"sync"
	"time"
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
		conn, err := server.listener.Accept()
		if err != nil {
			select {
			case <-server.quit:
				fmt.Print("problem with accepting connection", err)
				return
			default:
				fmt.Print("problem with accepting connection", err)
			}
		} else {
			server.wg.Add(1)
			go func() {
				server.handleConn(conn)
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

func (server *Server) handleConn(conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Errorf("Problem closing connection: %q", err)
		}
	}(conn)

ReadLoop:
	for {
		select {
		case <-server.quit:
			return
		default:
			conn.SetReadDeadline(time.Now().Add(200 * time.Millisecond))
			input, err := reader.Read(conn)
			if err != nil {
				if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
					continue ReadLoop
				}
				if err != io.EOF {
					fmt.Errorf("problem reading: %q", err)
					return
				}
			}
			commands, err := reader.ProcessData(input)

			if err != nil {
				fmt.Errorf("error occured while processing command: %v", err)
				continue ReadLoop
			}
			err = server.executeCommand(conn, commands)

			if err != nil {
				fmt.Errorf("error occured while executing command: %v", err)
				continue ReadLoop
			}
		}
	}
}

func (server *Server) executeCommand(conn net.Conn, commands [][]string) error {
	if len(commands) < 1 {
		return errors.New("Command should have at least one entry")
	}

	for i := 0; i < len(commands); i++ {
		if len(commands[i]) < 1 {
			return errors.New("Command should have at least one entry")
		}

		commandType := Command(strings.ToUpper(commands[i][0]))

		switch commandType {
		case Echo:
			conn.Write([]byte(protocol.WriteSimpleString("OK")))
		default:
			return fmt.Errorf("unsported command type: %v", commandType)
		}

	}

	return nil
}

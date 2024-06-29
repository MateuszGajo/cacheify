package server

import (
	"errors"
	"fmt"
	"io"
	"main/db"
	"main/protocol"
	"main/reader"
	"net"
	"os"
	"strconv"
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
	ECHO Command = "ECHO"
	SET  Command = "SET"
	GET  Command = "GET"
)

func CreateServer(address, port string) *Server {

	server := &Server{
		quit: make(chan interface{}),
	}

	listener, err := net.Listen("tcp", address+":"+port)
	fmt.Printf("\nServer running on %v:%v", address, port)
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

	db.CreateDb()

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
		commandArgs := commands[i][1:]

		switch commandType {
		case ECHO:
			server.handleEcho(conn, commandArgs)
		case SET:
			server.handleSet(conn, commandArgs)
		case GET:
			server.handleGet(conn, commandArgs)
		default:
			conn.Write([]byte(protocol.WriteSimpleError("Unknown command: " + string(commandType))))
		}

	}

	return nil
}

func (server *Server) handleEcho(conn net.Conn, args []string) {
	msg := args[0]
	conn.Write([]byte(protocol.WriteBulkString(msg)))
}

type ExpType string

const (
	PX ExpType = "px"
	EX ExpType = "ex"
)

func (server *Server) handleSet(conn net.Conn, args []string) {
	if len(args) < 2 {
		conn.Write([]byte(protocol.WriteSimpleError("We need at least key and value in set command")))
		return
	}
	key := args[0]
	value := args[1]
	if len(args) == 3 {
		conn.Write([]byte(protocol.WriteSimpleError("You need to add expiration value")))
		return
	}
	expMs := -1
	if len(args) == 4 {
		expType := ExpType(args[2])
		expValue, err := strconv.Atoi(args[3])
		if err != nil {
			conn.Write([]byte(protocol.WriteSimpleError("Invalid expiry time")))
			return
		}

		if expValue < 0 {
			conn.Write([]byte(protocol.WriteSimpleError("Expiry time need to be greater equal 0")))
			return
		}
		switch expType {
		case PX:
			expMs = expValue
		case EX:
			expMs = expValue * 1000
		default:
			conn.Write([]byte(protocol.WriteSimpleError(fmt.Sprintf("Exp type:%v is not supported", expType))))
		}
	}
	db.Set(key, value, expMs)
	conn.Write([]byte(protocol.WriteSimpleString("OK")))
}

func (server *Server) handleGet(conn net.Conn, args []string) {
	val, ok := db.Get(args[0])

	if !ok {
		conn.Write([]byte(protocol.WriteBulkString("-1")))
		return
	}
	conn.Write([]byte(protocol.WriteBulkString(val)))
}

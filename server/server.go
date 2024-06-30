package server

import (
	"errors"
	"fmt"
	"main/config"
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

type ServerRole string

const (
	REPLICA ServerRole = "REPLIC"
	MASTER  ServerRole = "MASTER"
)

type Server struct {
	listener      net.Listener
	quit          chan interface{}
	wg            sync.WaitGroup
	port          string
	address       string
	role          ServerRole
	masterAddress string
	masterPort    string
	replicaId     string
	replicaOffset string
}

type Command string

const (
	ECHO     Command = "ECHO"
	SET      Command = "SET"
	GET      Command = "GET"
	PING     Command = "PING"
	REPLCONF Command = "REPLCONF"
	PSYNC    Command = "PSYNC"
)

type ServerOptions func(*Server)

func WithPort(port string) ServerOptions {
	return func(s *Server) {
		s.port = port
	}
}

func WithAddress(address string) ServerOptions {
	return func(s *Server) {
		s.address = address
	}
}

func WithRole(role ServerRole) ServerOptions {
	return func(s *Server) {
		s.role = role
	}
}

func WithMaster(masterAddress, masterPort string) ServerOptions {
	return func(s *Server) {
		s.masterAddress = masterAddress
		s.masterPort = masterPort
	}
}

func CreateServer(options ...ServerOptions) *Server {

	server := &Server{
		quit:          make(chan interface{}),
		replicaId:     "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb",
		replicaOffset: "0",
	}

	for _, option := range options {
		option(server)
	}

	listener, err := net.Listen("tcp", server.address+":"+server.port)
	fmt.Printf("\nServer running on %v:%v", server.address, server.port)
	if err != nil {
		fmt.Print("can't run server", err)
		os.Exit(1)
	}

	server.listener = listener
	server.wg.Add(1)

	return server
}

func (server *Server) replicaHandshake() {
	conn, err := net.Dial("tcp", server.masterAddress+":"+server.masterPort)

	if err != nil {
		fmt.Printf("cant connect to replica, %v:%v because:%v", server.masterAddress, server.masterPort, err)
		os.Exit(1)
	}

	server.replicaHanshakePing(conn)
	server.replicaHanshakeRepl(conn)
	server.replicaHanshakePsync(conn)

	go func() {
		server.handleConn(conn)
	}()
}

func (server *Server) replicaHanshakeRepl(conn net.Conn) {
	_, err := conn.Write([]byte(protocol.WriteArrayString([]string{"REPLCONF", "listening-port", server.port})))

	if err != nil {
		fmt.Printf("cant send repl message to master: %v", err)
		os.Exit(1)
	}

	buf := make([]byte, 1024)

	n, err := conn.Read(buf)

	if err != nil {
		fmt.Printf("cant read response of replconf message: %v", err)
		os.Exit(1)
	}
	msg := buf[:n]

	if string(msg) != protocol.WriteSimpleString("OK") {
		fmt.Printf("Response to ping should be OK insted we got: %q", msg)
		os.Exit(1)
	}

	_, err = conn.Write([]byte(protocol.WriteArrayString([]string{"REPLCONF", "capa"})))

	if err != nil {
		fmt.Printf("cant send repl message to master: %v", err)
		os.Exit(1)
	}

	buf = make([]byte, 1024)

	n, err = conn.Read(buf)

	if err != nil {
		fmt.Printf("cant read response of replconf message: %v", err)
		os.Exit(1)
	}
	msg = buf[:n]

	if string(msg) != protocol.WriteSimpleString("OK") {
		fmt.Printf("Response to ping should be OK insted we got: %q", msg)
		os.Exit(1)
	}
}

func (server *Server) replicaHanshakePing(conn net.Conn) {
	_, err := conn.Write([]byte(protocol.WriteArrayString([]string{"PING"})))

	if err != nil {
		fmt.Printf("cant send echo message to master: %v", err)
		os.Exit(1)
	}

	buf := make([]byte, 1024)

	n, err := conn.Read(buf)

	if err != nil {
		fmt.Printf("cant read response of echo message: %v", err)
		os.Exit(1)
	}
	msg := buf[:n]

	if string(msg) != protocol.WriteSimpleString("PONG") {
		fmt.Printf("Response to ping should be pong insted we got: %q", msg)
		os.Exit(1)
	}
}
func (server *Server) replicaHanshakePsync(conn net.Conn) {
	_, err := conn.Write([]byte(protocol.WriteArrayString([]string{"PSYNC", "?", "-1"})))

	if err != nil {
		fmt.Printf("cant send echo message to master: %v", err)
		os.Exit(1)
	}

	buf := make([]byte, 1024)

	n, err := conn.Read(buf)

	if err != nil {
		fmt.Printf("cant read response of echo message: %v", err)
		os.Exit(1)
	}
	msg := strings.Split(string(buf[:n]), " ")
	fmt.Print(msg)
	respType := msg[0]
	replId := msg[1]

	if respType != "+FULLRESYNC" {
		fmt.Printf("expected response of psync to be type of +FULLRESYNC insted we got: %v", respType)
		os.Exit(1)
	}

	if len(replId) != 40 {
		fmt.Printf("Length of replica id should be 40 isnted we got: %v", len(replId))
		os.Exit(1)
	}

}

func (server *Server) RunServer() {
	defer server.wg.Done()

	db.CreateDb()

	if server.role == REPLICA {
		server.replicaHandshake()
		return
	}

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
				fmt.Printf("problem reading: %q", err)
				return
			}

			fmt.Println("input")
			fmt.Println(input)
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
		case PING:
			server.handlePing(conn, commandArgs)
		case SET:
			server.handleSet(conn, commandArgs)
		case GET:
			server.handleGet(conn, commandArgs)
		case REPLCONF:
			server.handleReplConf(conn, commandArgs)
		case PSYNC:
			server.handlePsync(conn, commandArgs)
		default:
			conn.Write([]byte(protocol.WriteSimpleError("Unknown command: " + string(commandType))))
		}

	}

	return nil
}

func (server *Server) handleEcho(conn net.Conn, args []string) {
	if len(args) == 0 {
		conn.Write([]byte(protocol.WriteSimpleError("We need passed value to echo command")))
		return
	}
	msg := args[0]
	conn.Write([]byte(protocol.WriteBulkString(msg)))
}

type ExpType string

const (
	PX ExpType = "px"
	EX ExpType = "ex"
)

type ReplType string

const (
	LISTENING_PORT ReplType = "LISTENING-PORT"
	CAPA           ReplType = "CAPA"
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

func (server *Server) handlePing(conn net.Conn, args []string) {
	conn.Write([]byte(protocol.WriteSimpleString("PONG")))
}

func (server *Server) handleReplConf(conn net.Conn, args []string) {

	replType := ReplType(strings.ToUpper(args[0]))

	switch replType {
	case LISTENING_PORT:
	case CAPA:
	default:
		conn.Write([]byte(protocol.WriteSimpleError("Doesnt support replType: " + string(replType))))
	}

	conn.Write([]byte(protocol.WriteSimpleString("OK")))

}

func (server *Server) handlePsync(conn net.Conn, args []string) {
	conn.Write([]byte(protocol.WriteSimpleString(fmt.Sprintf("%v %v %v", "FULLRESYNC", server.replicaId, "0"+config.CLRF))))
}

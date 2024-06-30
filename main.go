package main

import (
	"flag"
	"main/server"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

type ServerArgs struct {
	port          *string
	role          server.ServerRole
	masterAddress string
	masterPort    string
}

func getArgs() ServerArgs {
	port := flag.String("port", "6380", "Port for cache service")
	replicaOf := flag.String("replicaof", "", "Replica address")

	flag.Parse()

	role := server.MASTER

	if *replicaOf != "" {
		role = server.REPLICA
	}

	masterPort, masterAddress := "", ""
	if *replicaOf != "" {
		parts := strings.Split(*replicaOf, " ")
		if len(parts) == 2 {
			masterAddress = parts[0]
			masterPort = parts[1]
		}
	}

	return ServerArgs{
		port:          port,
		role:          role,
		masterAddress: masterAddress,
		masterPort:    masterPort,
	}
}

func main() {
	done := make(chan os.Signal, 1)
	args := getArgs()
	server := server.CreateServer(
		server.WithAddress("127.0.0.1"),
		server.WithPort(*args.port),
		server.WithRole(args.role),
		server.WithMaster(args.masterAddress, args.masterPort),
	)
	go func() {
		server.RunServer()
	}()

	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	<-done
	server.Close()
}

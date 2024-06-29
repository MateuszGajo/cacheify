package main

import (
	"flag"
	"fmt"
	"main/server"
	"os"
	"os/signal"
	"syscall"
)

type ServerArgs struct {
	port *string
}

func getArgs() ServerArgs {
	port := flag.String("port", "6380", "Port for cache service")

	flag.Parse()

	return ServerArgs{
		port: port,
	}
}

func main() {
	done := make(chan os.Signal, 1)
	args := getArgs()
	server := server.CreateServer("127.0.0.1", *args.port)
	go func() {
		server.RunServer()
	}()

	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	<-done
	fmt.Print("call done")
	server.Close()
}

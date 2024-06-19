package main

import (
	"main/server"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	done := make(chan os.Signal, 1)
	server := server.CreateServer("127.0.0.1", "6380")
	go func() {
		server.RunServer()
	}()

	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done

	server.Close()
}

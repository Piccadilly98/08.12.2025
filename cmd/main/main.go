package main

import (
	"log"

	"github.com/Piccadilly98/linksChecker/internal/server"
)

func main() {
	server := server.MakeServer(10)
	pid := server.Start("localhost:8080")
	log.Printf("Server work in PID: %d\n", pid)
	log.Printf("To shut down the server, run 'kill %d'\n", pid)
	<-server.ExitChan()
	log.Fatal("server stop")
}

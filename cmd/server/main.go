package main

import (
	"log"
	"net"

	"github.com/imperatrice00/oculis/internal"
	"github.com/imperatrice00/oculis/internal/handlers/server"
)

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:4000")
	if err != nil {
		log.Fatalln(err)
	}

	state := internal.NewState()

	conn, err := listener.Accept()
	if err != nil {
		log.Fatalln(err)
	}

	for {
		if err := server.HandlePacket(conn, state); err != nil {
			log.Println(err)
		}
	}
}

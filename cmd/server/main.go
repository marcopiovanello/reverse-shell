package main

import (
	"context"
	"flag"
	"log"

	"github.com/imperatrice00/oculis/internal/server"
)

func main() {
	addr := flag.String("c", "localhost:4000", "server address")
	flag.Parse()

	srv, err := server.NewClearTextServer(*addr)
	if err != nil {
		log.Fatalln(err)
	}

	srv.ReadLoop(context.Background())
}

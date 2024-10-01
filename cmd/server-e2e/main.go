package main

import (
	"context"
	"crypto/ecdh"
	"crypto/rand"
	"flag"
	"log"

	"github.com/imperatrice00/oculis/internal/server"
)

func main() {
	addr := flag.String("c", "localhost:4000", "server address")
	flag.Parse()

	curve := ecdh.P256()

	privKey, err := curve.GenerateKey(rand.Reader)
	if err != nil {
		log.Fatalln(err)
	}

	srv, err := server.NewE2EServer(*addr, privKey)
	if err != nil {
		log.Fatalln(err)
	}

	srv.ReadLoop(context.Background())
}

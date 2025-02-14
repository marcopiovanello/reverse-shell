package main

import (
	"context"
	"crypto/ecdh"
	"crypto/rand"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/imperatrice00/oculis/internal/server"
)

func main() {
	addr := flag.String("c", "localhost:4000", "server address")
	flag.Parse()

	curve := ecdh.X25519()

	privKey, err := curve.GenerateKey(rand.Reader)
	if err != nil {
		log.Fatalln(err)
	}

tryConn:
	srv, err := server.NewE2EServer(*addr, privKey)
	if err != nil {
		log.Println(err)
		log.Println("retrying...")
		time.Sleep(time.Second * 5)
		goto tryConn
	}

	log.Println("connected to", *addr)

	go gracefulShutdown(srv)

	srv.ReadLoop()
}

func gracefulShutdown(s server.Server) {
	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	go func() {
		<-ctx.Done()
		defer stop()

		s.Shutdown(ctx)
	}()
}

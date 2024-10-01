package server

import (
	"bytes"
	"context"
	"crypto/ecdh"
	"io"
	"log"
	"net"

	"github.com/imperatrice00/oculis/internal"
	"github.com/imperatrice00/oculis/internal/handlers"
)

type E2EServer struct {
	conn    net.Conn
	pubKey  *ecdh.PublicKey
	privKey *ecdh.PrivateKey
}

func NewE2EServer(addr string, privKey *ecdh.PrivateKey) (Server, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &E2EServer{
		conn:    conn,
		pubKey:  privKey.PublicKey(),
		privKey: privKey,
	}, nil
}

func (s *E2EServer) performKeyExchange() ([]byte, error) {
	// Prepare a buffer to store the server public key
	var buf bytes.Buffer

	// Copy the 65 byte of the public key
	_, err := io.CopyN(&buf, s.conn, internal.PUBLIC_KEY_SIZE)
	if err != nil {
		return nil, err
	}

	// Compute the curve and return the original ecdh public key (server)
	pubKey, err := ecdh.X25519().NewPublicKey(buf.Bytes())
	if err != nil {
		return nil, err
	}

	// Avoid leaking
	buf.Reset()

	// Compute the shared secret
	secret, err := s.privKey.ECDH(pubKey)
	if err != nil {
		return nil, err
	}

	// Exchange the client public key to the server
	s.conn.Write(s.privKey.PublicKey().Bytes())

	return secret, nil
}

func (s *E2EServer) ReadLoop(ctx context.Context) {
	state := internal.NewState()
	for {
		key, err := s.performKeyExchange()
		if err != nil {
			log.Println(err)
			return
		}
		if err := handlers.HandlePacket(s.conn, state, key); err != nil {
			log.Println(err)
			return
		}
	}
}

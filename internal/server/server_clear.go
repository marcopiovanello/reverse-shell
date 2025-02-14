package server

import (
	"context"
	"log"
	"net"

	"github.com/imperatrice00/oculis/internal"
	"github.com/imperatrice00/oculis/internal/handlers"
)

type ClearTextServer struct {
	conn net.Conn
}

func NewClearTextServer(addr string) (Server, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &ClearTextServer{
		conn: conn,
	}, nil
}

func (s *ClearTextServer) ReadLoop() {
	state := internal.NewState()
	for {
		if err := handlers.HandlePacket(s.conn, state, nil); err != nil {
			log.Println(err)
			return
		}
	}
}

func (s *ClearTextServer) Shutdown(ctx context.Context) {
	s.conn.Close()
}

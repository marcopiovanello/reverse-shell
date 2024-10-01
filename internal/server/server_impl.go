package server

import (
	"context"
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

func (s *ClearTextServer) ReadLoop(ctx context.Context) {
	state := internal.NewState()
	for {
		if err := handlers.HandlePacket(s.conn, state); err != nil {
			return
		}
	}
}

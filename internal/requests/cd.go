package requests

import (
	"net"

	"github.com/imperatrice00/oculis/internal"
	"github.com/imperatrice00/oculis/internal/client"
	"github.com/imperatrice00/oculis/internal/command"
)

func HandleChangeDirectory(path string) client.ClientHandlerFunc {
	return func(conn net.Conn, key []byte) error {
		req, err := internal.NewPacket(command.CD, []byte(path))
		if err != nil {
			return err
		}

		req.Write(conn)
		return nil
	}
}

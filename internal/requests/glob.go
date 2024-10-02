package requests

import (
	"bufio"
	"net"

	"github.com/imperatrice00/oculis/internal"
	"github.com/imperatrice00/oculis/internal/client"
	"github.com/imperatrice00/oculis/internal/command"
)

func HandleComputeGlob(path string, state *internal.ClientState) client.ClientHandlerFunc {
	return func(conn net.Conn, key []byte) error {
		req, err := internal.NewPacket(command.GLOB, []byte(path))
		if err != nil {
			return err
		}
		req.Write(conn)

		globJson, err := bufio.NewReader(conn).ReadBytes(internal.DELIMITER_SEQ)
		if err != nil {
			return err
		}

		state.Store("glob", globJson)

		return nil
	}
}

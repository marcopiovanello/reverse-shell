package requests

import (
	"bufio"
	"net"

	"github.com/imperatrice00/oculis/internal"
	"github.com/imperatrice00/oculis/internal/client"
	"github.com/imperatrice00/oculis/internal/command"
)

func HandleChangeDirectory(path string, state *internal.ClientState) client.ClientHandlerFunc {
	return func(conn net.Conn, key []byte) error {
		req, err := internal.NewPacket(command.CD, []byte(path))
		if err != nil {
			return err
		}

		req.Write(conn)

		pwdJson, err := bufio.NewReader(conn).ReadBytes(internal.DELIMITER_SEQ)
		if err != nil {
			return err
		}

		state.Store("pwd", pwdJson)

		return nil
	}
}

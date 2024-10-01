package requests

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"

	"github.com/imperatrice00/oculis/internal"
	"github.com/imperatrice00/oculis/internal/client"
	"github.com/imperatrice00/oculis/internal/command"
)

func HandleListDirectory(path string) client.ClientHandlerFunc {
	return func(conn net.Conn, key []byte) error {
		req, err := internal.NewPacket(command.LS, []byte(path))
		if err != nil {
			return err
		}

		req.Write(conn)

		res, err := bufio.NewReader(conn).ReadBytes(internal.DELIMITER_SEQ)
		if err != nil {
			return err
		}

		var files []string
		json.Unmarshal(res, &files)

		for _, file := range files {
			fmt.Println(file)
		}

		return nil
	}
}

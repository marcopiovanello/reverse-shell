package requests

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"

	"github.com/imperatrice00/oculis/internal"
	"github.com/imperatrice00/oculis/internal/client"
	"github.com/imperatrice00/oculis/internal/command"
)

func HandleListDirectory(path string) client.ClientHandlerFunc {
	return func(conn net.Conn) error {
		req, err := internal.NewPacket(command.LS, []byte(path))
		if err != nil {
			log.Println(err)
		}

		req.Write(conn)

		res, err := bufio.NewReader(conn).ReadBytes(internal.DELIMITER_SEQ)
		if err != nil {
			log.Println(err)
		}

		var files []string
		json.Unmarshal(res, &files)

		for _, file := range files {
			fmt.Println(file)
		}

		return nil
	}
}

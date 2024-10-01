package requests

import (
	"bufio"
	"encoding/json"
	"net"
	"path/filepath"

	"github.com/imperatrice00/oculis/internal"
	"github.com/imperatrice00/oculis/internal/client"
	"github.com/imperatrice00/oculis/internal/command"
)

func HandleGlobDownload(path, outputPath string) client.ClientHandlerFunc {
	return func(conn net.Conn, key []byte) error {
		req, err := internal.NewPacket(command.PWD, []byte{})
		if err != nil {
			return err
		}
		req.Write(conn)

		res, err := bufio.NewReader(conn).ReadBytes(internal.DELIMITER_SEQ)
		if err != nil {
			return err
		}

		var wd string
		if err := json.Unmarshal(res, &wd); err != nil {
			return err
		}

		files, err := filepath.Glob(filepath.Join(wd, path))
		if err != nil {
			return err
		}

		for _, file := range files {
			if key != nil {
				err := requestSingleAES(file, conn, outputPath, key)
				if err != nil {
					return err
				}
				return nil
			}
			err := requestSingle(file, conn, outputPath)
			if err != nil {
				return err
			}
		}

		return nil
	}
}

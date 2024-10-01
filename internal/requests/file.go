package requests

import (
	"encoding/binary"
	"errors"
	"io"
	"net"
	"os"
	"path/filepath"

	"github.com/imperatrice00/oculis/internal"
	"github.com/imperatrice00/oculis/internal/client"
	"github.com/imperatrice00/oculis/internal/command"
	"github.com/schollz/progressbar/v3"
)

func HandleFileDownload(path, outputPath string) client.ClientHandlerFunc {
	return func(conn net.Conn) error {
		return requestSingle(path, conn, outputPath)
	}
}

func requestSingle(path string, conn net.Conn, outputPath string) error {
	req, err := internal.NewPacket(command.DOWNLOAD, []byte(path))
	if err != nil {
		return err
	}

	req.Write(conn)

	var (
		chunkSize int32
		fileSize  int64
	)

	fd, err := os.Create(filepath.Join(outputPath, filepath.Base(path)))
	if err != nil {
		return err
	}

	binary.Read(conn, binary.LittleEndian, &fileSize)
	if fileSize == int64(-1) {
		return errors.New("an error occurred")
	}

	bar := progressbar.DefaultBytes(fileSize)
	writer := io.MultiWriter(fd, bar)

	for {
		binary.Read(conn, binary.LittleEndian, &chunkSize)
		_, err := io.CopyN(writer, conn, int64(chunkSize))
		if err != nil {
			break
		}
		if chunkSize < internal.CHUNK_SIZE {
			break
		}
	}

	return nil
}

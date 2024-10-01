package requests

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
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
	return func(conn net.Conn, key []byte) error {
		if key != nil {
			return requestSingleAES(path, conn, outputPath, key)
		}
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

func requestSingleAES(path string, conn net.Conn, outputPath string, key []byte) error {
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

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	nonceSize := aesGCM.NonceSize()

	for {
		buffer := &bytes.Buffer{}

		binary.Read(conn, binary.LittleEndian, &chunkSize)

		_, err := io.CopyN(buffer, conn, int64(chunkSize))
		if err != nil {
			break
		}

		nonce, ciphertext := buffer.Bytes()[:nonceSize], buffer.Bytes()[nonceSize:]

		plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
		if err != nil {
			return err
		}

		io.Copy(writer, bytes.NewReader(plaintext))

		if chunkSize < internal.CHUNK_SIZE {
			break
		}
	}

	return nil
}

package handlers

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/imperatrice00/oculis/internal"
)

func handleGetCurrentWorkingDirectory(conn net.Conn, state *internal.State) error {
	err := json.NewEncoder(conn).Encode(state.GetCurrentDirectory())
	if err != nil {
		return err
	}

	conn.Write(internal.DELIMITER_CONN)
	return nil
}

func handleChangeDirectory(conn net.Conn, payload []byte, state *internal.State) error {
	path := string(payload)

	if !filepath.IsAbs(path) || path == "" {
		path = filepath.Join(state.GetCurrentDirectory(), path)
	}

	stat, err := os.Stat(path)
	if err != nil {
		return err
	}

	if !stat.IsDir() {
		return errors.New(path + " is not a directory")
	}

	state.SetCurrentDirectory(filepath.Clean(path))

	json.NewEncoder(conn).Encode(state.GetCurrentDirectory())

	return nil
}

func handleListDirectory(conn net.Conn, payload []byte, state *internal.State) error {
	path := strings.TrimSpace(string(payload))

	if !filepath.IsAbs(path) || path == "" {
		path = filepath.Join(state.GetCurrentDirectory(), path)
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	files := make([]string, len(entries))

	for i, entry := range entries {
		files[i] = entry.Name()
	}

	json.NewEncoder(conn).Encode(files)
	conn.Write(internal.DELIMITER_CONN)

	return nil
}

/*
Request packet:

	2 bytes                    512 bytes
+---------+----------------------------------------------+
| command |                   payload                    |
+---------+----------------------------------------------+
|   cmd   | 256 bytes (filename) | 256 bytes (filesize)  |
+--------------------------------------------------------+

Response:
File chunked every 512KB
*/

func handleFileDownload(conn net.Conn, path []byte, state *internal.State) error {
	file := filepath.Clean(string(path))

	if !filepath.IsAbs(file) {
		file = filepath.Join(state.GetCurrentDirectory(), file)
	}

	fd, err := os.Open(file)
	if err != nil {
		handleFileError(conn)
		return err
	}
	stat, err := os.Stat(file)
	if err != nil {
		handleFileError(conn)
		return err
	}

	if stat.IsDir() {
		handleFileError(conn)
		return errors.New(file + " is not a file")
	}

	binary.Write(conn, binary.LittleEndian, stat.Size())

	br := bufio.NewReader(fd)
	buf := make([]byte, internal.CHUNK_SIZE)

	for {
		read, err := br.Read(buf)
		if err != nil {
			break
		}

		binary.Write(conn, binary.LittleEndian, int32(read))
		conn.Write(buf[:read])
	}

	return nil
}

func handleFileDownloadAES(conn net.Conn, path []byte, state *internal.State, key []byte) error {
	file := filepath.Clean(string(path))

	if !filepath.IsAbs(file) {
		file = filepath.Join(state.GetCurrentDirectory(), file)
	}

	fd, err := os.Open(file)
	if err != nil {
		handleFileError(conn)
		return err
	}
	stat, err := os.Stat(file)
	if err != nil {
		handleFileError(conn)
		return err
	}

	if stat.IsDir() {
		handleFileError(conn)
		return errors.New(file + " is not a file")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	binary.Write(conn, binary.LittleEndian, stat.Size())

	br := bufio.NewReader(fd)
	buf := make([]byte, internal.CHUNK_SIZE)

	for {
		read, err := br.Read(buf)
		if err != nil {
			break
		}

		nonce := make([]byte, aesGCM.NonceSize())
		if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
			return err
		}

		encryptedRead := aesGCM.Seal(nonce, nonce, buf[:read], nil)

		binary.Write(conn, binary.LittleEndian, int32(len(encryptedRead)))
		conn.Write(encryptedRead)
	}

	return nil
}

func handleFileError(conn net.Conn) {
	binary.Write(conn, binary.LittleEndian, int64(-1))
}

func handleGlobGeneration(conn net.Conn, glob []byte, state *internal.State) error {
	files, err := filepath.Glob(filepath.Join(state.GetCurrentDirectory(), string(glob)))
	if err != nil {
		log.Fatalln(err)
	}

	json.NewEncoder(conn).Encode(files)
	conn.Write(internal.DELIMITER_CONN)

	return nil
}

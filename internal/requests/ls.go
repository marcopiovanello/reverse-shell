package requests

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/fatih/color"
	"github.com/imperatrice00/oculis/internal"
	"github.com/imperatrice00/oculis/internal/client"
	"github.com/imperatrice00/oculis/internal/command"
	"github.com/imperatrice00/oculis/internal/responses"
)

// TODO: adapt to new response type
func HandleListDirectoryAES(path string) client.ClientHandlerFunc {
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

		block, err := aes.NewCipher(key)
		if err != nil {
			panic(err.Error())
		}

		aesGCM, err := cipher.NewGCM(block)
		if err != nil {
			panic(err.Error())
		}

		nonceSize := aesGCM.NonceSize()

		nonce, ciphertext := res[:nonceSize], res[nonceSize:]

		plaintextJson, err := aesGCM.Open(nil, nonce, ciphertext, nil)
		if err != nil {
			return err
		}

		var files []string
		json.Unmarshal(plaintextJson, &files)

		for _, file := range files {
			fmt.Println(file)
		}

		return nil
	}
}

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

		ls := &responses.DirectoryList{}
		json.Unmarshal(res, ls)

		fmt.Println(ls.Current)
		fmt.Println("total", len(ls.List))

		for _, file := range ls.List {
			filename := color.BlueString(file.Name)
			if file.IsDir {
				filename = color.GreenString(file.Name)
			}
			fmt.Printf("%d KiB\t%s %s\n",
				file.Size/1_000,
				file.MTime.Format(time.RFC822),
				filename,
			)
		}

		return nil
	}
}

package main

import (
	"bufio"
	"crypto/ecdh"
	"crypto/rand"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/imperatrice00/oculis/internal/client"
	"github.com/imperatrice00/oculis/internal/requests"
)

func main() {
	output := flag.String("o", "downloads", "where downloaded files will be saved")
	addr := flag.String("c", "localhost:4000", "server address")
	flag.Parse()

	curve := ecdh.P256()

	privKey, err := curve.GenerateKey(rand.Reader)
	if err != nil {
		log.Fatalln(err)
	}

	c, err := client.NewE2EClient(*addr, privKey, func(conn net.Conn) {
		log.Println(conn.RemoteAddr(), "has connected!")
	})
	if err != nil {
		log.Fatalln(err)
	}

	go c.Recoverer(func(c net.Conn) {
		log.Println("Peer", c.RemoteAddr(), "has connected!")
	})

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")

	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, "ls") {
			err := c.Send(requests.HandleListDirectory(""))
			if err != nil {
				log.Fatalln(err)
			}
		} else if strings.Contains(line, "cd") {
			args := strings.Split(line, "cd")
			if len(args) != 2 {
				continue
			}

			folder := strings.TrimSpace(args[1])
			err := c.Send(requests.HandleChangeDirectory(folder))
			if err != nil {
				log.Fatalln(err)
			}
		} else if strings.Contains(line, "download") {
			args := strings.Split(line, "download")
			if len(args) != 2 {
				continue
			}

			file := strings.TrimSpace(args[1])

			if strings.Contains(file, "*") {
				err := c.Send(requests.HandleGlobDownload(file, *output))
				if err != nil {
					log.Fatalln(err)
				}
			} else {
				err := c.Send(requests.HandleFileDownload(file, *output))
				if err != nil {
					log.Fatalln(err)
				}
			}
		} else {
			log.Println(scanner.Text(), "command not found")
		}

		fmt.Print("> ")
	}
}

package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/imperatrice00/oculis/internal/client"
	"github.com/imperatrice00/oculis/internal/requests"
)

func main() {
	output := flag.String("o", "downloads", "where downloaded files will be saved")
	addr := flag.String("c", "localhost:4000", "server address")
	flag.Parse()

	c, err := client.NewClearTextClient(*addr)
	if err != nil {
		log.Fatalln(err)
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")

	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, "ls") {
			c.Send(requests.HandleListDirectory(""))
		}
		if strings.Contains(line, "cd") {
			args := strings.Split(line, "cd")
			if len(args) != 2 {
				continue
			}

			folder := strings.TrimSpace(args[1])
			c.Send(requests.HandleChangeDirectory(folder))
		}
		if strings.Contains(line, "download") {
			args := strings.Split(line, "download")
			if len(args) != 2 {
				continue
			}

			file := strings.TrimSpace(args[1])

			if strings.Contains(file, "*") {
				c.Send(requests.HandleGlobDownload(file, *output))
			} else {
				c.Send(requests.HandleFileDownload(file, *output))
			}
		}

		fmt.Print("> ")
	}
}

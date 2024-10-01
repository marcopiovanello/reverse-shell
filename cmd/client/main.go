package main

import (
	"log"
	"net"

	"github.com/imperatrice00/oculis/internal/client"
	"github.com/imperatrice00/oculis/internal/requests"
)

func main() {
	c, err := client.NewClearTextClient("localhost:4000", func(conn net.Conn) {
		log.Println(conn.RemoteAddr(), "has connected!")
	})
	if err != nil {
		log.Fatalln(err)
	}

	c.Send(requests.HandleListDirectory("."))
	c.Send(requests.HandleListDirectory("."))
	c.Send(requests.HandleListDirectory("."))
	c.Send(requests.HandleFileDownload("yuzu.mp4", "downloads"))
}

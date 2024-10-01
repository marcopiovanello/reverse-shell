package client

import "net"

type ClientHandlerFunc = func(conn net.Conn) error

type Client interface {
	Send(hand ClientHandlerFunc) error
	Recoverer(onDiscovered func(net.Conn))
}

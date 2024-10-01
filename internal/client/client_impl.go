package client

import "net"

type ClearTextClient struct {
	conn net.Conn
}

func NewClearTextClient(addr string, onConnect func(net.Conn)) (Client, error) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	conn, err := listener.Accept()
	if err != nil {
		return nil, err
	}

	onConnect(conn)

	return &ClearTextClient{
		conn: conn,
	}, nil
}

func (c *ClearTextClient) Send(hand ClientHandlerFunc) error {
	return hand(c.conn)
}

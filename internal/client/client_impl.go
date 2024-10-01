package client

import "net"

type ClearTextClient struct {
	conn net.Conn
}

func NewClearTextClient(addr string) (Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &ClearTextClient{
		conn: conn,
	}, nil
}

func (c *ClearTextClient) Send(hand ClientHandlerFunc) error {
	return hand(c.conn)
}

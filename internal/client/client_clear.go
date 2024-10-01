package client

import "net"

type ClearTextClient struct {
	conn     net.Conn
	listener net.Listener
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
		conn:     conn,
		listener: listener,
	}, nil
}

func (c *ClearTextClient) Send(hand ClientHandlerFunc) error {
	return hand(c.conn, nil)
}

func (c *ClearTextClient) Recoverer(onDiscovered func(net.Conn)) {
	for {
		conn, err := c.listener.Accept()
		if err != nil {
			return
		}

		onDiscovered(conn)
		c.conn = conn
	}
}

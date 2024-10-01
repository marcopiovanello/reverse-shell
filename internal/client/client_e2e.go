package client

import (
	"bytes"
	"crypto/ecdh"
	"io"
	"net"

	"github.com/imperatrice00/oculis/internal"
)

type E2EClient struct {
	conn     net.Conn
	listener net.Listener
	pubKey   *ecdh.PublicKey
	privKey  *ecdh.PrivateKey
}

func NewE2EClient(addr string, privKey *ecdh.PrivateKey, onConnect func(net.Conn)) (Client, error) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	conn, err := listener.Accept()
	if err != nil {
		return nil, err
	}

	onConnect(conn)

	return &E2EClient{
		conn:     conn,
		listener: listener,
		pubKey:   privKey.PublicKey(),
		privKey:  privKey,
	}, nil
}

func (c *E2EClient) keyExchange(conn net.Conn) ([]byte, error) {
	var buf bytes.Buffer

	if _, err := conn.Write(c.pubKey.Bytes()); err != nil {
		return nil, err
	}

	io.CopyN(&buf, conn, internal.PUBLIC_KEY_SIZE)

	serverPubKey, err := ecdh.P256().NewPublicKey(buf.Bytes())
	if err != nil {
		return nil, err
	}

	secret, err := c.privKey.ECDH(serverPubKey)
	if err != nil {
		return nil, err
	}

	return secret, err
}

func (c *E2EClient) Send(hand ClientHandlerFunc) error {
	key, err := c.keyExchange(c.conn)
	if err != nil {
		return err
	}

	return hand(c.conn, key)
}

func (c *E2EClient) Recoverer(onDiscovered func(net.Conn)) {
	for {
		conn, err := c.listener.Accept()
		if err != nil {
			return
		}

		onDiscovered(conn)
		c.conn = conn
	}
}

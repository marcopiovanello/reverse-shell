package internal

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"

	"github.com/imperatrice00/oculis/internal/command"
)

const PAYLOAD_SIZE = 256

type Packet struct {
	Command command.Command
	Payload [PAYLOAD_SIZE]byte
}

func NewPacket(cmd command.Command, payload []byte) (*Packet, error) {
	if len(payload) > PAYLOAD_SIZE {
		return nil, errors.New("payload too large")
	}

	container := [PAYLOAD_SIZE]byte{}
	copy(container[:], payload)

	return &Packet{
		Command: cmd,
		Payload: container,
	}, nil
}

func (p *Packet) CleanPayload() []byte {
	cut := bytes.IndexByte(p.Payload[:], byte(rune(0)))
	return p.Payload[0:cut]
}

func (p *Packet) Write(w io.Writer) {
	binary.Write(w, binary.LittleEndian, p)
}

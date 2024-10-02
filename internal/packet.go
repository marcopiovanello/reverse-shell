package internal

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"

	"github.com/imperatrice00/oculis/internal/command"
)

const PAYLOAD_SIZE = 512

type Packet struct {
	Cmd command.Command
	Pyl [PAYLOAD_SIZE]byte
}

func NewPacket(cmd command.Command, payload []byte) (*Packet, error) {
	if len(payload) > PAYLOAD_SIZE {
		return nil, errors.New("payload too large")
	}

	container := [PAYLOAD_SIZE]byte{}
	copy(container[:], payload)

	return &Packet{
		Cmd: cmd,
		Pyl: container,
	}, nil
}

func (p *Packet) Payload() (path []byte, extra []byte) {
	path = p.Pyl[:255]
	extra = p.Pyl[256:]

	cut := bytes.IndexByte(path, byte(rune(0)))
	path = path[0:cut]

	cut = bytes.IndexByte(extra, byte(rune(0)))
	extra = extra[0:cut]

	return
}

func (p *Packet) Command() command.Command {
	return p.Cmd
}

func (p *Packet) Write(w io.Writer) {
	binary.Write(w, binary.LittleEndian, p)
}

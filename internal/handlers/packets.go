package handlers

import (
	"encoding/binary"
	"errors"
	"net"

	"github.com/imperatrice00/oculis/internal"
	"github.com/imperatrice00/oculis/internal/command"
)

func HandlePacket(conn net.Conn, state *internal.State, secret []byte) error {
	packet := &internal.Packet{}

	err := binary.Read(conn, binary.LittleEndian, packet)
	if err != nil {
		return err
	}

	switch packet.Command {
	case command.CD:
		payload := packet.CleanPayload()
		return handleChangeDirectory(conn, payload, state)
	case command.LS:
		payload := packet.CleanPayload()
		return handleListDirectory(conn, payload, state)
	case command.DOWNLOAD:
		payload := packet.CleanPayload()
		if secret != nil {
			return handleFileDownloadAES(conn, payload, state, secret)
		}
		return handleFileDownload(conn, payload, state)
	case command.PWD:
		return handleGetCurrentWorkingDirectory(conn, state)
	case command.DOWNLOAD_DIR:
		// TODO: implement handler
		return errors.New("unimplemented method")
	case command.QUIT:
		// TODO: implement handler
		return errors.New("unimplemented method")
	default:
		return errors.New("unimplemented method")
	}
}

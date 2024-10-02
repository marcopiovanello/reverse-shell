package command

type Command = uint16

const (
	LS Command = iota + 1
	CD
	PWD
	GLOB
	DOWNLOAD
	DOWNLOAD_DIR
	QUIT
)

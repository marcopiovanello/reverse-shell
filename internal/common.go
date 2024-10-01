package internal

const (
	CHUNK_SIZE      int32 = 512 * 1000
	DELIMITER_SEQ   byte  = '\n'
	PUBLIC_KEY_SIZE       = 32 // X25519 -> 32 | P256 -> (32*2)+1
)

var DELIMITER_CONN = []byte("\n")

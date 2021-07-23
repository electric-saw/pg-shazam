package state

import (
	"encoding/binary"
)

func IntToByteArray(u uint32) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint32(buf, u)
	return buf
}

func ByteArrayToInt(b []byte) uint32 {
	return binary.BigEndian.Uint32(b)
}

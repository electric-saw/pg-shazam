package state

import (
	"fmt"
)

type HashSet struct {
	Database string
	Table    string
	Fields   []string
}

func (h *HashSet) Encode() []byte {
	var buf []byte

	buf = append(buf, 'H')

	buf = append(buf, IntToByteArray(uint32(len(h.Database)))...)
	buf = append(buf, []byte(h.Database)...)
	buf = append(buf, IntToByteArray(uint32(len(h.Table)))...)
	buf = append(buf, []byte(h.Table)...)
	buf = append(buf, IntToByteArray(uint32(len(h.Fields)))...)
	for _, field := range h.Fields {
		buf = append(buf, IntToByteArray(uint32(len(field)))...)
		buf = append(buf, []byte(field)...)
	}
	return buf
}

func (h *HashSet) Decode(buf []byte) error {
	var idx = uint32(0)

	if buf[idx] != 'H' {
		return fmt.Errorf("Invalid buffer type %v", buf[idx])
	}

	idx++

	var tmpSize uint32

	tmpSize = ByteArrayToInt(buf[idx : idx+8])
	idx += 8

	h.Database = string(buf[idx : idx+tmpSize])
	idx += tmpSize

	tmpSize = ByteArrayToInt(buf[idx : idx+8])
	idx += 8

	h.Table = string(buf[idx : idx+tmpSize])
	idx += tmpSize

	fieldCount := ByteArrayToInt(buf[idx : idx+8])
	idx += 8

	for i := uint32(0); i < fieldCount; i++ {
		tmpSize = ByteArrayToInt(buf[idx : idx+8])
		idx += 8

		field := string(buf[idx : idx+tmpSize])
		idx += tmpSize
		h.Fields = append(h.Fields, field)
	}
	return nil
}

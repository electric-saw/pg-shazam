package state

import (
	"fmt"
)

type Session struct {
	NodeId string
	PID    uint32
	Secret uint32
}

func (s *Session) Encode() []byte {
	var buf []byte

	buf = append(buf, 'S')

	buf = append(buf, IntToByteArray(uint32(len(s.NodeId)))...)
	buf = append(buf, []byte(s.NodeId)...)
	buf = append(buf, IntToByteArray(s.PID)...)
	buf = append(buf, IntToByteArray(uint32(s.Secret))...)
	return buf
}

func (s *Session) Decode(buf []byte) error {
	var idx = uint32(0)

	if buf[idx] != 'S' {
		return fmt.Errorf("Invalid buffer type %s", string(buf[idx]))
	}

	idx++

	tmpSize := ByteArrayToInt(buf[idx : idx+8])
	idx += 8

	s.NodeId = string(buf[idx : idx+tmpSize])
	idx += tmpSize

	s.PID = ByteArrayToInt(buf[idx : idx+8])
	idx += 8

	s.Secret = uint32(ByteArrayToInt(buf[idx : idx+8]))
	return nil
}

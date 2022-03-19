package protocol

import (
	"encoding/binary"
	"errors"
)

type stream struct {
}

func (s *stream) MsgLen(data []byte) (length uint64, err error) {
	if len(data) <= 4 {
		err = errors.New("head len short")
		return
	}

	l := binary.BigEndian.Uint32(data[0:4])
	if len(data) < int(4+l) {
		err = errors.New("msg short")
		return
	}

	length += uint64(l) + 4
	return
}

func (s *stream) Unpack(data []byte) (msg []byte, err error) {
	msg = data[4:]
	return
}

func (s *stream) Pack(data []byte) (msg []byte, err error) {
	var head []byte
	binary.BigEndian.PutUint32(head, uint32(len(data)))
	msg = append(head, data...)
	return
}

func NewStream() *stream {
	return &stream{}
}

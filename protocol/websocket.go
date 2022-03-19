package protocol

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"github/alwaysLinger/seco/ctx"
)

const (
	OpcodeInvalid = 0x00
	OpcodeText    = 0x01
	OpcodeBin     = 0x02
	OpcodeClose   = 0x08
	OpcodePing    = 0x09
	OpcodePong    = 0x0A
)

type Websocket struct {
	*Http
	handShaked bool
	// isFin      bool
	opCode    uint8
	headerLen uint8
}

func (w *Websocket) reset() {
	w.opCode = OpcodeInvalid
	w.headerLen = 0
}

func (w *Websocket) GetOpCode() uint8 {
	return w.opCode
}

func (w *Websocket) MsgLen(data []byte) (length uint64, err error) {
	if !w.handShaked {
		return w.Http.MsgLen(data)
	} else {
		w.reset()
		dataLen := uint64(len(data))
		if dataLen <= 6 {
			err = errors.New("websocket header short")
			return
		}
		finFlag := data[0] >> 7 & 0x01
		if finFlag == 0 {
			err = errors.New("not fin frame")
			return
		}
		w.opCode = data[0] & 0b00001111
		if w.opCode == OpcodeClose {
			err = errors.New("close frame")
			return
		}
		isMask := data[1] >> 7 & 0x01
		if isMask == 0 {
			err = errors.New("client frame not masked")
			return
		}
		pLen := data[1] & 0b01111111
		if pLen <= 125 {
			length = uint64(pLen) + 6
			w.headerLen = 6
		}
		if dataLen < length {
			err = errors.New("not enough frame data")
			return
		}
		if pLen == 126 {
			Len := binary.BigEndian.Uint16(data[2:4])
			length = uint64(Len) + 8
			w.headerLen = 8
		}
		if dataLen < length {
			err = errors.New("not enough frame data")
			return
		}
		if pLen >= 127 {
			Len := binary.BigEndian.Uint64(data[2:10])
			length = Len + 14
			w.headerLen = 14
		}
		if dataLen < length {
			err = errors.New("not enough frame data")
			return
		}
	}
	return
}

func (w *Websocket) Unpack(data []byte) (msg []byte, err error) {
	if !w.handShaked {
		return w.Http.Unpack(data)
	}
	defer w.reset()
	mask := data[w.headerLen-4 : w.headerLen]
	payload := data[w.headerLen:]
	for i, v := range payload {
		payload[i] = mask[i&0x03] ^ v
	}
	msg = payload
	return
}

func (w *Websocket) Pack(data []byte) (msg []byte, err error) {
	dataLen := len(data)
	if dataLen <= 125 {
		frame := make([]byte, dataLen+2)
		frame[0] = 0b10000000 | OpcodeText
		frame[1] = byte(dataLen)
		frame = append(frame[0:2], data...)
		msg = frame
		return
	} else if dataLen <= 0xFFFF {
		frame := make([]byte, dataLen+4)
		frame[0] = 0b10000000 | OpcodeText
		frame[1] = 126
		binary.BigEndian.PutUint16(frame[2:4], uint16(dataLen))
		frame = append(frame[0:4], data...)
		msg = frame
		return
	} else {
		frame := make([]byte, dataLen+10)
		frame[0] = 0b10000000 | OpcodeText
		frame[1] = 127
		binary.BigEndian.PutUint64(frame[2:10], uint64(dataLen))
		frame = append(frame[0:10], data...)
		msg = frame
		return
	}
}

func (w *Websocket) HasHandShaked() bool {
	return w.handShaked
}

func (w *Websocket) HandeShake(req *ctx.Request) string {
	key := req.GetHeader("Sec-WebSocket-Key")
	if key == "" {
		return ""
	}
	acceptSumArg := key + "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
	acceptKeySha := sha1.Sum([]byte(acceptSumArg))
	acceptKeyStr := base64.StdEncoding.EncodeToString(acceptKeySha[:])
	return "HTTP/1.1 101 Switching Protocols\r\nUpgrade: websocket\r\nConnection: Upgrade\r\nSec-WebSocket-Accept: " + acceptKeyStr +
		"\r\nSec-WebSocket-Version: 13\r\nServer: seco\r\n\r\n"
}

func (w *Websocket) SetHandeShake() {
	w.handShaked = true
}

func NewWebsocket() *Websocket {
	return &Websocket{}
}

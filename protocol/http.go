package protocol

import (
	"bytes"
	"errors"
	"github.com/alwaysLinger/seco/ctx"
	"strconv"
)

type HttpError uint8

const (
	HttpNoData HttpError = iota + 1
	HttpMsgLongError
	HttpMsgNotfoundSplit
	HttpMsgNotfinish
)

type Http struct {
	err HttpError
}

func (h *Http) MsgLen(data []byte) (length uint64, err error) {
	dataLen := len(data)
	if dataLen == 0 {
		h.err = HttpNoData
		err = errors.New("no http data")
		return
	}
	bodySplitIdx := bytes.Index(data, []byte("\r\n\r\n"))
	if bodySplitIdx <= 0 {
		h.err = HttpMsgNotfoundSplit
		err = errors.New("\r\n\r\n not found")
		return
	}

	if dataLen > 1024*1024 {
		h.err = HttpMsgLongError
		err = errors.New("http data too long")
		return
	}

	// get content-length
	cntIdx := bytes.Index(data, []byte("Content-Length"))
	if cntIdx == -1 {
		length = uint64(bodySplitIdx) + 4
		// fmt.Println("length:", length)
		return
	}

	cntData := data[cntIdx:]
	firstN := bytes.IndexByte(cntData, byte('\n'))
	cntBytes := cntData[16 : firstN-1]
	cntLen, _ := strconv.Atoi(string(cntBytes))
	length = uint64(bodySplitIdx) + 4 + uint64(cntLen)
	if length > uint64(dataLen) {
		h.err = HttpMsgNotfinish
		err = errors.New("http msg not enough")
		return
	}

	return
}

func (h *Http) Unpack(data []byte) (msg []byte, err error) {
	msg = data
	return
}

func (h *Http) Pack(data []byte) (smg []byte, err error) {
	return
}

func (h *Http) Request(data []byte) *ctx.Request {
	return ctx.NewRequest(data)
}

func (h *Http) GetErr() HttpError {
	return h.err
}

func NewHttp() *Http {
	return &Http{}
}

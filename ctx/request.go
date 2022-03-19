package ctx

import (
	"bytes"
)

type Request struct {
	method   string
	uri      string
	version  string
	headers  map[string]string
	segments map[string]string
	fields   map[string]string
}

func (r *Request) praseHeaders() map[string]string {
	return nil
}

func (r *Request) parseSegments() map[string]string {
	return nil
}

func (r *Request) parseFields() map[string]string {
	return nil
}

func (r *Request) GetHeader(header string) string {
	if v, ok := r.headers[header]; ok {
		return v
	}
	return ""
}

func (r *Request) IsKeepAlive() bool {
	// close
	// if "Keep-Alive" == r.GetHeader("Connection") {
	// 	return true
	// }
	// return false
	return true
}

// NewRequest simply form-urlencoded implement
func NewRequest(raw []byte) *Request {
	r := &Request{
		headers:  make(map[string]string, 1),
		segments: make(map[string]string, 1),
		fields:   make(map[string]string, 1),
	}
	bodyIdx := bytes.Index(raw, []byte("\r\n\r\n"))
	headerCnt := raw[0:bodyIdx]
	bodyCnt := raw[bodyIdx+4:]
	headers := bytes.Split(headerCnt, []byte("\r\n"))
	for i, v := range headers {
		if i == 0 {
			uriIdx := bytes.IndexByte(v, byte(' '))
			r.method = string(v[0:uriIdx])
			versionIdx := bytes.LastIndexByte(v, byte(' '))
			r.version = string(v[versionIdx+1:])
		} else {
			hkv := bytes.Split(v, []byte(": "))
			r.headers[string(hkv[0])] = string(hkv[1])
		}
	}

	if len(bodyCnt) != 0 {

	}

	return r
}

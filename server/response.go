package server

import (
	"strconv"
)

type Response struct {
	conn *Connection
}

func (r *Response) buildResponse(resp string) []byte {
	resp = "HTTP/1.1 200 OK\r\nServer: beco\r\nContent-Length: " + strconv.Itoa(len(resp)) + "\r\n\r\n" + resp
	return []byte(resp)
}

func (r *Response) Send(resp string) (length int, err error) {
	return r.conn.Send(r.buildResponse(resp))
}

func (r *Response) HandeShake(key string) bool {
	if _, err := r.conn.Send([]byte(key)); err != nil {
		return false
	}
	return true
}

func NewResponse(c *Connection) *Response {
	return &Response{
		conn: c,
	}
}

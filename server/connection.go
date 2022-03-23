package server

import (
	"fmt"
	"github.com/alwaysLinger/seco/protocol"
	"io"
	"net"
)

type Connection struct {
	conn    net.Conn
	addr    string
	server  *Server
	protocl string
	packer  protocol.Protocol
	running bool
	buffer  [1024 * 1024]byte
	nlast   uint64
}

func (c *Connection) GetAddr() string {
	return c.addr
}

func (c *Connection) HandlePayload() {
	for c.running {
		l, err := c.conn.Read(c.buffer[c.nlast:])
		if err != nil {
			if err == io.EOF {
				c.close()
				return
			}
		}
		c.nlast += uint64(l)
		if c.protocl == "tcp4" {
			c.server.Call("receiveCb", c, c.buffer[0:c.nlast])
			c.nlast = 0
		} else {
			c.ProtoclHandle()
		}
	}
}

func (c *Connection) ProtoclHandle() {
	tmp := c.buffer[0:c.nlast]
	for c.nlast > 0 {
		length, err := c.packer.MsgLen(tmp)
		// break
		if err != nil {
			fmt.Println("msglen error:", err)
			// 如果是http协议 那么直接断开连接
			if c.protocl == "http" {
				p := c.packer.(*protocol.Http)
				switch p.GetErr() {
				case protocol.HttpMsgLongError:
					c.close()
					break
				case protocol.HttpNoData, protocol.HttpMsgNotfoundSplit, protocol.HttpMsgNotfinish:
					break
				}
				// break
			}
			if c.protocl == "ws" {
				if c.packer.(*protocol.Websocket).GetOpCode() == protocol.OpcodeClose {
					c.close()
					break
				}
			}
			break
		}
		if c.protocl == "tcp4" {
			data, _ := c.packer.Unpack(tmp[0:length])
			c.server.Call("receiveCb", c, data)
		} else if c.protocl == "http" {
			data, _ := c.packer.Unpack(tmp[0:length])
			req := c.packer.(*protocol.Http).Request(data)
			resp := NewResponse(c)
			c.server.Call("requestCb", req, resp)
			if req.IsKeepAlive() {
				c.close()
				break
			}
		} else {
			packer := c.packer.(*protocol.Websocket)
			data, _ := packer.Unpack(tmp[0:length])
			if !packer.HasHandShaked() {
				req := packer.Request(data)
				resp := NewResponse(c)
				if !resp.HandeShake(packer.HandeShake(req)) {
					c.close()
					break
				}
				packer.SetHandeShake()
			} else {
				// fmt.Println("length:", length)
				// fmt.Println(string(data), len(data))
				c.server.Call("messageCb", c.server, c, data)
			}
		}

		tmp = tmp[length:]
		c.nlast -= length

		if c.nlast == 0 {
			break
		}
	}
}

func (c *Connection) Send(data []byte) (int, error) {

	return c.conn.Write(data)
}

func (c *Connection) close() {
	c.server.Close(c)
	_ = c.conn.Close()
	c.running = false
	// c.nlast = 0
}

func (c *Connection) Push(data []byte) (int, error) {
	frame, _ := c.packer.Pack(data)
	fmt.Println(frame)
	return c.conn.Write(frame)
}

func NewClient(conn net.Conn, s *Server) (c *Connection, err error) {
	c = &Connection{
		conn:    conn,
		addr:    conn.RemoteAddr().String(),
		server:  s,
		protocl: s.GetProtocl(),
		running: true,
	}
	switch c.protocl {
	case "stream":
		c.packer = protocol.NewStream()
	case "http":
		c.packer = protocol.NewHttp()
	case "ws":
		c.packer = protocol.NewWebsocket()
	}
	return
}

package server

import (
	"fmt"
	"github/alwaysLinger/seco/ctx"
	"net"
	"sync"
)

type Server struct {
	protocol     string
	addr         string
	conns        map[string]*Connection
	clientNum    int64
	maxClientNum int64
	l            net.Listener
	mux          sync.Mutex

	errCb   func(err string)
	startCb func(s *Server)

	connectCb func(s *Server, conn *Connection)
	receiveCb func(s *Server, conn *Connection, data []byte)
	closeCb   func(s *Server, conn *Connection)

	messageCb func(s *Server, conn *Connection, data []byte)
	openCb    func(s *Server, conn *Connection)

	requestCb func(req *ctx.Request, resp *Response)
}

func NewServer(protocol string, addr string, maxClient int64) *Server {
	return &Server{
		protocol:     protocol,
		addr:         addr,
		conns:        make(map[string]*Connection, 10),
		maxClientNum: maxClient,
	}
}

func (s *Server) OnConnect(cb func(s *Server, conn *Connection)) {
	s.connectCb = cb
}

func (s *Server) OnReceive(cb func(s *Server, conn *Connection, data []byte)) {
	s.receiveCb = cb
}

func (s *Server) OnClose(cb func(s *Server, conn *Connection)) {
	s.closeCb = cb
}

func (s *Server) OnRequest(cb func(req *ctx.Request, resp *Response)) {
	s.requestCb = cb
}

func (s *Server) OnMessage(cb func(s *Server, conn *Connection, data []byte)) {
	s.messageCb = cb
}

func (s *Server) Call(name string, args ...interface{}) {
	switch name {
	case "errCb":
		s.errCb(args[0].(string))
	case "connectCb":
		s.connectCb(s, args[0].(*Connection))
	case "receiveCb":
		s.receiveCb(s, args[0].(*Connection), args[1].([]byte))
	case "closeCb":
		s.closeCb(s, args[1].(*Connection))
	case "requestCb":
		s.requestCb(args[0].(*ctx.Request), args[1].(*Response))
	case "messageCb":
		s.messageCb(s, args[1].(*Connection), args[2].([]byte))
	}
}

func (s *Server) Start() {
	fmt.Printf("listen on %s\n", s.addr)

	l, err := net.Listen("tcp4", s.addr)
	if err != nil {
		s.Call("errCb", err.Error())
		return
	}
	s.l = l
	s.loop()
}

func (s *Server) GetProtocl() string {
	return s.protocol
}

func (s *Server) loop() {
	for {
		conn, err := s.l.Accept()
		if err != nil {
			s.Call("errCb", err.Error())
			return
		}
		c, err := NewClient(conn, s)
		if err != nil {
			s.Call("errCb", err.Error())
		}
		s.addConn(c)
		s.Call("connectCb", c)
		go c.HandlePayload()
	}
}

func (s *Server) addConn(c *Connection) {
	defer s.mux.Unlock()
	s.mux.Lock()
	s.conns[c.GetAddr()] = c
}

func (s *Server) Close(conn *Connection) {
	s.Call("closeCb", s, conn)
	defer s.mux.Unlock()
	s.mux.Lock()
	if _, ok := s.conns[conn.GetAddr()]; ok {
		delete(s.conns, conn.GetAddr())
	}
}

package main

import (
	"bytes"
	"fmt"
	"github/alwaysLinger/seco/ctx"
	"github/alwaysLinger/seco/server"
	"net"
)

func mainasdasd() {
	l, err := net.Listen("tcp4", ":9527")
	if err != nil {
		fmt.Println("listen failed")
	}

	// b := [5]byte{1,2,3,4}
	// b = [5]byte{1,2,3,4}
	fmt.Println(l)

	// buf:= bytes.Buffer{}
	var buf []byte
	_ = bytes.NewBuffer(buf)
	fmt.Printf("%[1]v %[1]T %d %d\n", buf, len(buf), cap(buf))
	// nb.

	bb := make([]byte, 10)
	bc := bb[5:7]
	fmt.Println(len(bb), len(bc))

	var abc []byte
	abc = make([]byte, 20)
	bbb := bytes.NewBuffer(abc[10:])
	bbb.Write([]byte("abc"))
	fmt.Println(bbb)
	// fmt.Println(abc)
	// ccc := bbb.Read
	// 数组通过截取的操作 得到的是对应类型的切片
	// 修改切片就是修改底层引用的数组

	var abcd [10]byte
	read(abcd[0:], 123)
	fmt.Println(abcd)
	read(abcd[3:], 123)
	fmt.Println(abcd)
}

// 这种性能还是非常高的
func read(aa []byte, x byte) {
	aa[0] = x
}

func main() {
	s := server.NewServer("ws", "127.0.0.1:9527", 1024)

	s.OnConnect(func(s *server.Server, conn *server.Connection) {
		// fmt.Println(conn)
	})

	// s.OnReceive(func(s *server.Server, conn *server.Connection, data []byte) {
	// 	fmt.Println(string(data))
	// })

	s.OnRequest(func(req *ctx.Request, resp *server.Response) {
		fmt.Println(req.GetHeader("Host"))
		resp.Send("nihao世界")
	})

	s.OnMessage(func(s *server.Server, conn *server.Connection, data []byte) {
		fmt.Print(string(data))
		conn.Push([]byte("nihao世界"))
	})

	s.OnClose(func(s *server.Server, conn *server.Connection) {
		fmt.Println(conn.GetAddr() + "--closed")
	})

	s.Start()
}

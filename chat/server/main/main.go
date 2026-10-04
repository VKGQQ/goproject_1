package main

import (
	"fmt"
	"net"
)

func handleConn(conn net.Conn) {
	defer conn.Close()
	processor := &Processor{
		Conn: conn,
	}
	err := processor.serve()
	if err != nil {
		fmt.Println("客户端与服务器通讯协程出错：", err)
		return
	}
}

func main() {
	fmt.Println("服务器在8889端口监听")
	listen, err := net.Listen("tcp", "0.0.0.0:8889")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer listen.Close()
	for {
		fmt.Println("等待客户端连接服务器......")
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println(err)
		}
		go handleConn(conn)
	}
}

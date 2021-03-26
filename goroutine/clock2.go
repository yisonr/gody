package main

import (
	"io"
	"log"
	"net"
	"time"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		// 监听器的 accept 方法被阻塞, 直到有连接请求进来
		if err != nil {
			log.Print(err) // 例: 连接中止
			continue
		}
		go handleConn(conn) // 一次处理一个连接
	}
}

func handleConn(c net.Conn) {
	defer c.Close() // 延迟的 close 调用关闭自己的连接
	for {
		_, err := io.WriteString(c, time.Now().Format("15:04:05\n"))
		if err != nil {
			return // 例: 连接断开
		}
		time.Sleep(1 * time.Second)
	}
}

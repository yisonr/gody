package main

import (
	"io"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	go mustCopy(os.Stdout, conn) // 第二个 goroutine
	mustCopy(conn, os.Stdin)     // 主 goroutine
	// 主 goroutine 从标准输入读取并发送到服务器后,
	// 第二个 goroutine 读取服务器的回复并输出
}

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}

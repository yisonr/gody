package main

// 回声服务器

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
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
		go handleConn(conn) // 处理多个客户端连接
	}
}

func echo(c net.Conn, shout string, delay time.Duration) {
	fmt.Fprintln(c, "\t", strings.ToUpper(shout))
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", shout)
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", strings.ToLower(shout))
}

func handleConn(c net.Conn) {
	input := bufio.NewScanner(c)
	for input.Scan() {
		go echo(c, input.Text(), 10*time.Second)
		// 回声并发
	}
	// note: ignore input.Error() error
	c.Close()
}

package main

// p178 练习8.3

import (
	// "fmt"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	tcpAddr, err := net.ResolveTCPAddr("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(tcpAddr.IP)
	// fmt.Println(tcpAddr.Port)
	// fmt.Println(tcpAddr.Zone)
	conn, err := net.DialTCP("tcp", &net.TCPAddr{}, tcpAddr)
	if err != nil {
		log.Fatal(err)
	}
	done := make(chan struct{})
	go func() {
		io.Copy(os.Stdout, conn) // 忽略错误处理
		log.Println("done")
		done <- struct{}{} // 指示主 goroutine
	}()
	mustCopy(conn, os.Stdin)
	// conn.Close()
	conn.CloseWrite()
	<-done
}

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}

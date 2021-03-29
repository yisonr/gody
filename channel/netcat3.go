package main

// 使用无缓冲通道进行的通信导致发送和接收 goroutine
// 同步化，无缓冲通道也称为同步通道

import (
	"fmt"
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
	done := make(chan struct{})
	// 有时候通信本身以及通信发生的时间也很重要，
	// 当强调这方面的时候，把消息叫做事件(event)
	// 当事件没有携带额外的信息时，它单纯的进行同步
	// 通过使用 struct{} 元素类型的通道来强调它
	go func() {
		io.Copy(os.Stdout, conn) // 忽略错误处理
		log.Println("done")
		// 程序总是在退出前记录 "done" 消息
		done <- struct{}{} // 指示主 goroutine
	}()
	mustCopy(conn, os.Stdin)
	conn.Close()
	<-done // 等待后台 goroutine 完成
	// 此通道为无缓冲通道
	// 接收阻塞，直到从 goroutine 发送值到通道
}

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}

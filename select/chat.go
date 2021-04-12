package main

// 聊天服务器
// TODO:
//  30行, 定义client 的意义
//  在 handleConn 函数中, ch 的接收和阻塞???
//  为什么有的地方接收发送符号在规范格式化后不靠近通道一侧

// todo: 理解
// 当有 n 个客户 session 在连接的时候，程序并发
// 运行着 2n+2 个相互通信的 goroutine, 它不需要
// 隐式的加锁操作。clients map 限制在广播器这一
// 个 goroutine 中被访问，所以不会并发访问它。
// 唯一被多个 goroutine 共享的变量是通道以及
// net.Conn 的实例，它们又都是并发安全的。

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	go broadcaster() // 广播器
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

type client chan<- string // 对外发送消息的通道

var (
	entering = make(chan client) // 只能发送???
	leaving  = make(chan client)
	messages = make(chan string) // 所有接收的客户消息
)

func broadcaster() {
	clients := make(map[client]bool) // 所有连接的客户端
	for {
		select {
		case msg := <-messages:
			// 把所有接收的消息广播给所有的客户
			// 发送消息通道
			for cli := range clients {
				cli <- msg // 检查工具不能判断其中哪个是通道???
			}
		case cli := <-entering:
			clients[cli] = true
		case cli := <-leaving:
			delete(clients, cli)
			close(cli)
		}
	}
}

func handleConn(conn net.Conn) {
	ch := make(chan string)   // 对外发送客户消息的通道
	go clientWriter(conn, ch) // 这个 goroutine 通过ch接收主goroutine发送的值

	who := conn.RemoteAddr().String()
	ch <- "You are " + who
	messages <- who + " has arrived" // 发送给所有连接的客户
	entering <- ch                   // 一个通道发送给另一个通道???

	input := bufio.NewScanner(conn)
	for input.Scan() {
		messages <- who + ": " + input.Text() // 广播某个客户端的消息
	}
	// 忽略了 input.Err() 中可能的错误

	leaving <- ch // ch 为什么不会阻塞???
	messages <- who + " has left"
	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg) // 忽略网络层面的错误
	}
}

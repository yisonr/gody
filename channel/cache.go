package main

import (
	// "bytes"
	"fmt"
	"log"
	"net/http"
)

func main() {
	addr := mirroredQuery()
	fmt.Println(addr)
}

func mirroredQuery() string {
	responses := make(chan string, 3)
	// 使用缓冲通道并发的向三个镜像地址发送请求，
	// 只接收到第一个响应后便退出goroutine
	go func() { responses <- request("http://www.baidu.com") }()
	go func() { responses <- request("https://www.bing.com") }()
	go func() { responses <- request("https://blog.csdn.net") }()
	// 如果使用一个无缓冲通道，两个比较慢的 goroutine 将被卡住
	// 因为他们发送响应结果到通道的时候没有 goroutine 来接收，
	// 称为 goroutine 泄露， 属于一个bug
	// 应确保 goroutine 在不需要的时候可以自动结束
	return <-responses
}

func request(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	if err != nil {
		log.Fatalln(err)
	}
	return resp.Status + " " + url
}

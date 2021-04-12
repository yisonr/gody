package main

// TODO: 添加更多的request, 添加goroutine的取消功能, p197 练习8.11

// 无缓冲通道和缓冲通道的选择、缓冲通道容量大小的选择，都会对程序
// 的正确性产生影响。
// 无缓冲通道提供强同步保障，因为每一次发送都需要和每一次对应的接收同步
// 对于缓冲通道，这些操作则是解耦的。在内存无法提供缓冲容量的情况下，
// 可能会导致程序死锁。

// 缓冲的通道也可能影响程序的性能。
// 对于 goroutine1-->goroutine2-->goroutine3
//  gr(goroutine)1 的速度相比 gr2 的速度快(非量纲)，他们之间的通道常满
//  gr1 的速度相比 gr2 的速度慢，他们之间的通道常空
// 这样 gr1 或 gr2 其中之一就是常态阻塞, 这样缓冲区的存在是没有价值的。
// 同时对于 gr1-->gr2-->gr3, 若 gr2 的速度相比其他两者都慢(gr2跟不上gr1的
// 供应，或者满足不了 gr3 的需要)，可以再加一个 gr4
// 帮助 gr2 这段流程，独立地执行相同的任务。

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

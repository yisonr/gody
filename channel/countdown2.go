package main

import (
	"fmt"
	"os"
	"time"
)

// 在10s倒计时中可以按回车键终止发送

func main() {
	abort := make(chan struct{})
	go func() {
		os.Stdin.Read(make([]byte, 1))
		abort <- struct{}{}
	}()

	fmt.Println("Commencing countdown. Press return to abort.")
	// select 多路复用
	// 两个通道的接收放到主goroutine中，只有都满足才能结束程序
	// 因为倒计时操作和回车操作都会导致阻塞, 需要使用多路复用检测
	// 最早接收的通道, 并执行相应的操作
	// select 一直等待，直到一次通信来告知有一些情况
	// 可以执行，然后进行这次通信，执行此情况对应的语句
	// 其他的通信将不会发生，对于没有对应情况的select，
	// select{} 将永远等待.(select也有default分支)
	// 此 select 语句等待两个事件中第一个到达的事件，中止事件或
	// 者指示事件过去10s的事件。
	select {
	case <-time.After(10 * time.Second):
		// After 函数立即返回一个通道，然后启动一个新的 goroutine
		// 在间隔指定时间后发送一个值到通道

		// ...不执行操作
	case <-abort:
		fmt.Println("Launch aborted!")
		return
	}
	fmt.Println("Start launch!")
	// launch()
}

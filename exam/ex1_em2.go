package main

// 面试题
// channel 是用于在 goroutine 之间同步和传递数据的关键数据类型
// 注: channel 在这里存在滥用

// csp: 通信顺序进程
// csp 是一种消息传递模型, 通过在 goroutine 之间传递数据
// 来传递消息，而不是对数据进行加锁来实现同步访问

import (
	"fmt"
	"runtime"
)

func main() {
	// 分配一个逻辑处理器给调度器使用
	runtime.GOMAXPROCS(1)

	numCh := make(chan int, 1)
	strCh := make(chan string, 1)
	numCh <- 1
	strCh <- "hello"

	select {
	case value := <-numCh:
		fmt.Println(value)
	case value := <-strCh:
		panic(value)
	}

	// select 在多个情况同时满足的情况下，会随机选择一个
}

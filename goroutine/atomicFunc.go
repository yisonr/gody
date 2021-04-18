package main

// Load 和 Store 类函数提供对数值类型的安全访问
// 程序使用 LoadInt64 和 StoreInt64 创建一个同步标志，
// 这个标志可以向程序里多个 goroutine 通知某个特殊状态

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

var (
	// 通知正在执行的 goroutine 停止工作的标志
	shutdown int64
	wg       sync.WaitGroup
)

func main() {
	wg.Add(2)

	go doWork("A")
	go doWork("B")

	// 给定 goroutine 执行的时间
	time.Sleep(1 * time.Second)

	// 该停止工作了， 安全的设置 shutdown 标志
	atomic.StoreInt64(&shutdown, 1)
	// 如果哪个 doWork goroutine 试图在 main 函数中
	// 调用 StoreInt64 的同时调用 LoadInt64 函数，那么
	// 原子函数会将这些调用互相同步，保证这些操作都是
	// 安全的，不会进入竞争状态

	wg.Wait()
}

// doWork 模拟执行工作的 goroutine,
// 检测 shutdown 标志来决定是否提前终止
func doWork(name string) {
	defer wg.Done()

	for {
		fmt.Printf("Doing %s Work\n", name)
		time.Sleep(250 * time.Millisecond)

		// 检测 shutdown 确定是否停止工作
		if atomic.LoadInt64(&shutdown) == 1 {
			fmt.Printf("Shutting %s Down\n", name)
			break
		}
	}
}

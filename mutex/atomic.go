package main

// 在 ./race.go 的基础上讨论同步问题:
// go 提供了传统的同步 goroutine 的进制，就是对
// 共享资源加锁
// 使用 atomic 包来提供对数值类型的安全访问
// 原子函数能够以很底层的加锁机制来同步访问
// 整形变量和指针

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
)

var (
	// goroutine 都要增加其值的变量
	counter int64

	// wg 等待 goroutine 完成
	wg sync.WaitGroup // 计数信号量
)

func main() {
	// 要等待两个 goroutine
	wg.Add(2)

	// 创建两个 goroutine
	go incCounter(1)
	go incCounter(2)

	// 等待 goroutine 结束
	wg.Wait()
	fmt.Println("Final Counter:", counter)
}

// incCounter 增加包里的 counter 变量的值
func incCounter(id int) {
	// 退出时调用 Done 通知 main 函数工作已经完成
	defer wg.Done()

	for count := 0; count < 2; count++ {
		// 安全的对 counter 的值加1
		// 注意: 这种只能同步访问整形变量和指针
		atomic.AddInt64(&counter, 1) // 原子函数
		// 该方法强制同一时刻只能有一个 goroutine 运行
		// 并完成这个加法操作
		// 当 goroutine 试图去调用任何原子函数时，这些
		// goroutine 都会自动根据所引用的变量做同步处理

		// 当前 goroutine 从线程退出，并放回到队列
		runtime.Gosched()
	}
}

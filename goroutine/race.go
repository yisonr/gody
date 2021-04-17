package main

// 展示如何造成竞争状态
// 应避免这种情况

import (
	"fmt"
	"runtime"
	"sync"
)

var (
	// goroutine 都要增加其值的变量
	counter int

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
		// 捕获 counter 的值
		value := counter

		// 当前 goroutine 从线程退出，并放回到队列
		runtime.Gosched()

		// 增加本地 value 变量的值
		value++

		// 写 counter
		counter = value
	}
}

package main

// 对一个共享资源的读和写必须是原子化的，换句话说，同一时刻只能有
// 一个 goroutine 对共享资源进行读和写操作

import (
	"fmt"
	"runtime"
	"sync"
)

var (
	// 所有 goroutine 都要增加其值的变量
	counter int

	// wg 用来等待程序结束
	wg sync.WaitGroup
)

func main() {
	wg.Add(2)

	go inCounter(1)
	go inCounter(2)

	wg.Wait()
	fmt.Println("Final Counter:", counter)
}

func inCounter(id int) {
	defer wg.Done()

	for count := 0; count < 2; count++ {
		// 捕获 counter 的值
		value := counter

		// 从当前 goroutine 从线程退出，并放回到队列
		runtime.Gosched()

		// 增加本地value变量的值
		value++

		// 将该值保存回 counter
		counter = value
	}
}

// runtime 包的 Gosched 函数，用于将 goroutine 从当前线程退出，给其他
// goroutine 运行的机会, 在两次操作中间这样做强制调度器切换两个  goroutine
// 以便让竞争状态的效果更加明显.

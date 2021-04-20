package main

// 使用互斥锁来定义一段需要同步访问的
// 代码临界区资源的同步访问

import (
	"fmt"
	"runtime"
	"sync"
)

var (
	// counter 是所有 goroutine 都要增加其值的变量
	counter int
	wg      sync.WaitGroup
	// mutex 用来定义一段代码临界区
	mutex sync.Mutex
)

func main() {
	wg.Add(2)

	go incCounter(1) // 1 号 goroutine
	go incCounter(2) // 2 号 goroutine

	wg.Wait()
	fmt.Printf("Final Counter: %d\n", counter)
}

// incCounter 使用互斥锁来同步并保证安全访问
// 增加包里 counter 变量的值
func incCounter(id int) {
	// id 只为区分不同的 goroutine 而用，无实际意义
	defer wg.Done()

	for count := 0; count < 2; count++ {
		// 同一时刻只允许一个 goroutine 进入此临界区
		mutex.Lock() // 加锁
		{            // 临界区代码，此处的括号非必需
			// 捕获 counter 的值
			value := counter
			// 当前 goroutine 从线程退出，并放回到队列
			runtime.Gosched()
			// 强制将当前 goroutine 退出当前线程后，
			// 调度器会再次分配这个 goroutine 继续运行
			// 当程序结束时，得到正确的值:4, 竞争状态不存在
			value++
			counter = value
		}
		mutex.Unlock()
		// 释放锁，允许其他等待的 goroutine 进入临界区
	}
}

package main

// 展示如何造成竞争状态
// 应避免这种情况

// go build -race ./race.go
// 用竞态检测器标志来编译程序
// 运行编译后的程序会展示引发数据竞争的
// goroutine 以及相关代码

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
		value := counter // 对 counter 变量的读操作， 存在数据竞态

		// 当前 goroutine 从线程退出，并放回到队列
		runtime.Gosched()
		/*
			GoSched 只是让出处理器来允许其他 goroutine 运行，
			并没有挂起此 goroutine, 所以此 goroutine 将会自动
			的恢复执行
		*/

		// Gosched 函数用于将 goroutine 从当前线程退出，
		// 给其他 goroutine 运行的机会, 使用此函数的目的
		// 是强调调度器切换两个 goroutine,
		// 以展示竞争状态的效果

		// 切换 goroutine, 再次切换回来后从这里开始运行
		// 每个 goroutine 都会覆盖另一个 goroutine 的工作,
		// 这种覆盖发生在 goroutine 切换的时候，
		// 每个 goroutine 创造一个 counter 变量的副本，
		// 之后切换到另一个 goroutine . 当此 goroutine
		// 再次运行的时候， counter 变量的值已经变了，
		// 但是 goroutine 并没有更新自己的那个副本的值，
		// 而是继续用这个副本的值进行递增后来更新 counter
		// 变量，结果覆盖了另一个 goroutine 的工作

		// 增加本地 value 变量的值
		value++

		// 写 counter
		counter = value // 对 counter 变量的写操作, 存在数据竞态
	}
}

// 变量 counter 会进行4次读写操作，每个 goroutine 执行两次
// 执行结果: 2

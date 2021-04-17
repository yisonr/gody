package main

// 并发
// 调度器的行为
// 调度器在操作系统之上，将操作系统的线程与语言
// 运行时的逻辑处理器绑定，并在逻辑处理器上允许 goroutine
// 调度器在任何给定的时间，都会全面控制哪个 goroutine 要在哪个
// 逻辑处理器上运行.

// 进程是一个包含了应用程序在运行中需要用到和维护的各种资源的容器
// 这些资源包括但不限于内存地址空间、文件和设备的句柄以及线程.
// 一个线程是一个执行空间，这个空间会被操作系统调度来执行函数中
// 所写的代码。
// 每个进程至少包含一个线程，每个进程的初始化线程被称为主线程。
// 操作系统将线程调度到某个处理器上运行，这个处理器并不一定是进程
// 所在的处理器.

// 并发: 在同一个逻辑处理器上，交替执行两个或多个事件，即在同一时间间隔发生
// 并行: 在不同的逻辑处理器上，同时执行两个或多个事件，即在同一时刻发生

// todo: 加深理解
// * 只有在有多个逻辑处理器且可以同时让每个 goroutine
// 运行在一个可用的物理处理器上的时候, goroutine 才会并行运行???

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	// 分配一个逻辑处理器给调度器使用 -- 并发
	runtime.GOMAXPROCS(1)
	// runtime.GOMAXPROCS(runtime.NumCPU()) // 8核, 并行
	// 给每个可用的核心(物理处理器)分配一个逻辑处理器, 让 goroutine 并行运行

	// wg 用来等待程序完成
	// 计数加 2, 表示要等待两个 goroutine
	var wg sync.WaitGroup // 计数信号量
	wg.Add(2)

	fmt.Println("Start Goroutines!")

	go func() {
		// 在函数退出时调用 done 来通知 main 函数工作已经完成
		defer wg.Done()

		// 显示字母表3次
		for count := 0; count < 3; count++ {
			for char := 'a'; char < 'a'+26; char++ {
				fmt.Printf("%c ", char)
			}
		}
	}()

	go func() {
		// 在函数退出时调用 done 来通知 main 函数工作已经完成
		defer wg.Done()

		// 显示字母表3次
		for count := 0; count < 3; count++ {
			for char := 'A'; char < 'A'+26; char++ {
				fmt.Printf("%c ", char)
			}
		}
	}()

	// todo:
	// 这两个 goroutine 进入调度队列的顺序?
	// 实验数据得出第二个 goroutine 先运行

	// 等待 goroutine 结束
	fmt.Println("Waiting to finish!")
	wg.Wait()

	fmt.Println("\nTerminating program")
}

package main

// todo: 调度器原理

// goroutine 在逻辑处理器上执行，而逻辑处理器
// 具有独立的系统线程和运行队列

// 展示 goroutine 调度器是如何在单个线程上
// 切分时间片的

// 并发
// go 语言运行时会把 goroutine 调度到逻辑处理器上
// 运行，这个逻辑处理器绑定到唯一的操作系统线程，
// 当 goroutine 可以运行的时候，会被放入逻辑处理器的
// 执行队列中。

// 如果想让 goroutine 并行，必须使用多于一个逻辑处理器，
// 当有多个逻辑处理器时，调度器会将 goroutine 平等分配到
// 每个逻辑处理器上, 这会让 goroutine 在不同的线程上运行

import (
	"fmt"
	"runtime"
	"sync"
)

var wg sync.WaitGroup // 等待程序完成

func main() {
	// 分配一个逻辑处理器给调度器使用
	runtime.GOMAXPROCS(1)
	wg.Add(2)

	// 创建两个 goroutine
	fmt.Println("Create goroutines")
	go printPrime("A")
	go printPrime("B")
	// 调度时间片
	// B goroutine 先运行，打印到素数 5399 左右开始切换
	// 到 A goroutine, 每次运行时，调度器切换的时间点
	// 会有些许不同

	fmt.Println("Waiting to finish")
	wg.Wait()

	fmt.Println("Terminating program")
}

// printPrime 显示 5000 以内的素数值
func printPrime(prefix string) {
	defer wg.Done()

next: // 循环标签
	for outer := 2; outer < 6000; outer++ {
		// todo:
		// 如何查看 goroutine 调度器的信息
		// 如何查看 goroutine 调度时间片?
		// 和进程的时间片调度的区别
		for inner := 2; inner < outer; inner++ {
			if outer%inner == 0 {
				continue next
			}
		}
		fmt.Printf("%s:%d\n", prefix, outer)
	}
	fmt.Println("Completed", prefix)
}

package main

// 调度器的行为

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	// 分配一个逻辑处理器给调度器使用 -- 并发
	runtime.GOMAXPROCS(1)

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

	// 等待 goroutine 结束
	fmt.Println("Waiting to finish!")
	wg.Wait()

	fmt.Println("\nTerminating program")
}

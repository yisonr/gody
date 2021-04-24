package main

/*
有缓冲通道
有缓冲通道是一种在接收前能存储一个或
多个值的通道，这种类型的通道并不强制
要求 goroutine 之间必须同时完成发送和
接收.
通道会阻塞发送和接收动作的条件也不同
*/

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

const (
	// 4个 goroutine 处理10个工作
	numberGoroutines = 4  // 要使用的 goroutine 的数量
	taskLoad         = 10 // 要处理的工作的数量
)

var wg sync.WaitGroup

// 包级别的初始化函数
func init() {
	// 初始化随机数种子
	rand.Seed(time.Now().Unix())
}

func main() {
	// 分配一个逻辑处理器给调度器
	// runtime.GOMAXPROCS(1)
	// 调度时间片，如果使用并发模式，可能第一个 goroutine 在
	// 调度时间片到达之前就已经完成了10个工作
	tasks := make(chan string, taskLoad)

	// 等待 goroutines
	wg.Add(numberGoroutines)
	for gr := 1; gr <= numberGoroutines; gr++ {
		go worker(tasks, gr)
	}

	// 增加一组要完成的工作
	for post := 1; post <= taskLoad; post++ {
		tasks <- fmt.Sprintf("Task: %d", post)
	}

	// 所有工作都完成时关闭通道
	// 以便所有 goroutine 退出
	// 发送结束后即关闭通道
	close(tasks)

	// 等待所有工作完成
	wg.Wait()
}

// worker 作为 goroutine 启动来处理
// 从有缓冲的通道传入的工作
func worker(tasks chan string, worker int) {
	// 通知函数已经返回
	defer wg.Done()

	for {
		// 等待分配工作
		task, ok := <-tasks
		if !ok {
			// 通道关闭
			fmt.Println(task + "------")
			fmt.Printf("Worker: %d: Shutting Down\n", worker)
			return
		}

		// 显示开始工作
		fmt.Printf("Worker: %d: Started %s\n", worker, task)

		// 随机等待一段时间模拟工作
		sleep := rand.Int63n(100)
		time.Sleep(time.Duration(sleep) * time.Millisecond)

		// 显示已经完成工作
		fmt.Printf("Worker: %d: Completed %s\n", worker, task)

	}
}

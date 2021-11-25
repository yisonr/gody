package main

// 控制程序的生命周期

import (
	"log"
	"os"
	"time"
)

// runner 作为包导出使用的情况下, 只导出了 Start 和 Add 方法,
// 避免了 run 和 gotInterrupt 被引用, TODO: 好处?

const timeout = 12 * time.Second

func main() {
	log.Println("Start work.")
	/*
	 * 在 main 中初始化通道并添加需要监控的任务,
	 * 并接收系统中断信号, 在单独的 goroutine 中执行任务,
	 * 在 main 中监控 goroutine 中任务的执行
	 * main 中发送系统中断信号给 goroutine, goroutine 发送任务运行
	 * 结果(中断/完成)给 main
	 * main 中监控程序的执行时间(r.timeout)
	 *
	 * r.complete 和 r.timeout 都需要接收到信号后立即在main中进行处理,
	 * 所以他们的 channel 缓冲区容量为0
	 *
	 */

	// 为本次执行分配超时时间
	r := New(timeout) // 3 秒

	// 添加要执行的任务
	r.Add(CreatTask(), CreatTask(), CreatTask())

	// 执行任务并处理结果
	if err := r.Start(); err != nil {
		// 任务没有顺利完成, 超时/中断
		switch err {
		case ErrInterrupt:
			log.Println("Terminating due to interrupt.")
			os.Exit(2)
		case ErrTimeout:
			log.Println("Terminating due to timeout.")
			os.Exit(1)
		}
	}
	log.Println("Process ended.")
}

func CreatTask() func(int) {
	return func(id int) {
		log.Printf("Processor - Task #%d.", id)
		time.Sleep(time.Duration(id) * time.Second)
	}
}

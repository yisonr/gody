package main

import (
	"log"
	"os"
	"time"
)

const timeout = 6 * time.Second

func main() {
	log.Println("Start work.")

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

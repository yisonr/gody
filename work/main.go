package main

// 展示使用 worker 创建一个  goroutine 池并完成工作

import (
	"fmt"
	"log"
	"sync"
	"time"
)

var names = []string{
	"steve",
	"bob",
	"mary",
	"theres",
	"jason",
}

type namePrinter struct {
	name string
}

func (n *namePrinter) Task() {
	log.Println(n.name)
	time.Sleep(time.Second) // 模拟任务的执行
}

func main() {
	startTime := time.Now()
	// maxGoroutine = 2
	p := New(2)

	var wg sync.WaitGroup
	// main 函数中用 wg 管理自己创建的 goroutine, 这些 goroutine 提交任务至
	// Pool 中, Pool 中使用自己的 wg 字段，管理自己创建的 goroutine, 这些
	// goroutine 从 Pool 的 work 通道中接收任务并执行
	wg.Add(10 * len(names))
	// 10 * len(names) 的 goroutine 竞争提交任务至池中, Pool中使用 2 个
	for i := 0; i < 10; i++ {
		for _, name := range names {
			// 创建一个 worker(namePrinter) 并提供指定的名字
			np := namePrinter{
				name: name,
			}
			go func() {
				// 将任务提交执行, Run 返回时表明任务已经处理完成,
				// 释放 goroutine
				p.Run(&np)
				wg.Done()
			}()
		}
	}
	wg.Wait() // 保证提交任务的 goroutine  完成
	fmt.Printf("%v, main goroutine finish\n", time.Since(startTime))

	// 让工作池停止工作, 等待现有的工作完成(保证执行任务的goroutine完成)
	p.Shutdown()
	fmt.Printf("%v, pool goroutine finish\n", time.Since(startTime))
}

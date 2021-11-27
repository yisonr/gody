package main

/*
 * 展示如何使用无缓冲的通道创建一个 goroutine 池, 这些 goroutine 执行并控制
 * 一组工作, 让其并发执行;
 *
 * 这种情况既不需要一个工作队列, 也不需要一组 goroutine 配合执行,
 * 无缓冲的通到保证两个 goroutine  之间的数据交换, 这种使用无缓冲
 * 通道的方法允许使用者知道何时 goroutine 池正在执行工作, 而且如
 * 果池里的所有 goroutine 都忙, 无法接受新的工作的时候, 也能及时
 * 通过通道来通知调用者, 使用无缓冲的通道不会有工作在队列里丢失或
 * 卡住, 所有工作都会被处理.
 *
 */

import "sync"

// Worker 必须满足接口类型, 才能使用工作池
type Worker interface {
	Task()
}

// Pool 提供一个 goroutine 池, 这个池可以完成任何已提交的 worker 任务
type Pool struct {
	work chan Worker
	wg   sync.WaitGroup
}

func New(maxGoroutine int) *Pool {
	p := Pool{
		work: make(chan Worker),
	}

	p.wg.Add(maxGoroutine)
	for i := 0; i < maxGoroutine; i++ {
		go func() {
			for w := range p.work {
				w.Task()
			}
			p.wg.Done()
		}()
	}
	return &p
}

// Run 提交工作到工作池
func (p *Pool) Run(w Worker) {
	p.work <- w
}

// 等待所有的 goroutine 停止工作
func (p *Pool) Shutdown() {
	close(p.work) // 通知 goroutine 发送完成
	p.wg.Wait()
}

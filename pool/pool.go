package main

import (
	"errors"
	"io"
	"log"
	"sync"
)

// Pool 管理一组可以安全地在多个
// goroutine 间共享的资源，被管理
// 的资源必须实现 io.Closer 接口
type Pool struct {
	m         sync.Mutex
	resources chan io.Closer            // interface
	factory   func() (io.Closer, error) // 为池创建新资源
	closed    bool
}

// ErrPoolClosed 表示请求(Acquire) 了一个
// 已经关闭的池
var ErrPoolClosed = errors.New("Pool has been closed.")

// New 创建一个用来管理资源的池
// 这个池需要一个可以分配新资源的函数,
// 并规定池的大小
// 池的工厂函数
func New(fn func() (io.Closer, error), size uint) (*Pool, error) {
	if size <= 0 {
		return nil, errors.New("Size value too small.")
	}

	return &Pool{
		factory:   fn,
		resources: make(chan io.Closer, size),
	}, nil
}

// Acquire 从池中获取一个资源
func (p *Pool) Acquire() (io.Closer, error) {
	select {
	// 检查是否有空闲的资源
	case r, ok := <-p.resources:
		log.Println("Acquire:", "Shared Resource")
		if !ok {
			// 通道关闭
			return nil, ErrPoolClosed
		}
		return r, nil
		// 因为没有空闲的资源可用，所以提供一个新资源
	default:
		log.Println("Acquire:", "New Resource")
		return p.factory()
	}
}

// Release 将一个使用后的资源放回池里
func (p *Pool) Release(r io.Closer) {
	// 保证本操作和 close 操作的安全
	p.m.Lock()
	defer p.m.Unlock()

	// 如果池已经被关闭，销毁这个资源
	if p.closed {
		r.Close()
		return
	}

	select {
	//  将资源放入队列
	case p.resources <- r:
		log.Println("Release:", "In Queue")

		//  如果队列已满，则关闭这个资源
	default:
		log.Println("Release:", "Closing")
		r.Close()
	}
}

// Close 会让资源池停止工作, 并关闭所有现有的资源
func (p *Pool) Close() {
	// 保证本操作与 Release 操作的安全
	// 在同一时刻，只能有一个 goroutine 执行这段代码
	p.m.Lock()
	defer p.m.Unlock()

	// 如果 pool 已被关闭，什么也不做
	if p.closed {
		return
	}

	// 将池关闭
	p.closed = true

	// 在清空通道里的资源之前，将通道关闭
	// 否则会死锁
	close(p.resources)
	// 关闭资源
	for r := range p.resources {
		r.Close()
	}
}

// Release 和 Close 方法中的互斥量是同一个互斥量，这样可以
// 阻止这两个方法在不同的 goroutine 里同时运行
// 1. 对 closed 的访问存在数据竞态
// 2. 应该保证不往一个已经关闭的通道里发送数据

package main

/*
 * 展示使用有缓冲的通道实现资源池, 来管理可以在任意数量的 goroutine 之间
 * 共享及独立使用的资源, 这种模式适用于需要共享一组静态资源的情况(如共享数据库
 * 连接或者内存缓冲区); 如果 goroutine 需要从池中得到这些资源中的一个, 必须从
 * 池里申请, 使用完后归还到资源池
 */

import (
	"errors"
	"io"
	"log"
	"sync"
)

/*
 * 读取关闭的无缓冲通道, 返回值为对应通道类型的零值和false
 * 读取关闭的有缓冲通道, 会读取通道中的值, 如果通道为空时, 返回值为
 * 对应通道类型的零值和false
 * 使用 range 遍历通道, 通道写完后必须关闭通道, 否则 range 遍历会锁死
 */

// Pool 管理一组可以安全地在多个 goroutine 间共享的资源，被管理
// 的资源必须实现 io.Closer 接口(很多资源, 如数据库连接或打开的
// 文件都是需要被 Close() 的)
type Pool struct {
	m         sync.Mutex // 对 closed 的读写加锁
	resources chan io.Closer
	// 资源通道, 由 main 中的 Close 负责关闭和清空
	// 由 goroutine 中的 Acquire 负责接收, Release 负责发送
	factory func() (io.Closer, error) // 为池创建新资源
	closed  bool                      // 决定了资源池的关闭
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
		// 若通道关闭, 但是通道内仍有资源, 则资源
		// 可以被接收; 直至通道为空时, ok 为 false
		// 若通道未关闭, 但通道内无可用资源, 则运行
		// default 分支
		log.Println("Acquire:", "Shared Resource")
		if !ok {
			// 通道关闭
			return nil, ErrPoolClosed
		}
		return r, nil
	default:
		log.Println("Acquire:", "New Resource")
		// 因为没有空闲的资源可用，所以提供一个新资源
		// 这个新资源被用完后, 会归还到资源池里, 若资源
		// 池已被关闭或资源池已满, 则销毁这个资源
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

		//  如果队列已满，则销毁这个资源
	default:
		log.Println("Release:", "Closing")
		r.Close()
	}
}

// Close 会让资源池停止工作, 并关闭所有现有的资源
// 在主 main 中完成任务后关闭池, 与归还资源的操作
// 存在对 p.closed 读写的数据竞态, 所以在 Close 和
// Release 使用同一个互斥量对操作进行加锁

// 否则可能会存在 goroutine(调用 Release ) 往一个已经关闭的通道
// (p.resources)里发送数据, 从而导致奔溃(line 84)
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

	// 为了使 range 遍历不死锁, 必须关闭通道, 同时防止资源池不断的提供新资源
	// 做无效的操作后销毁, 继续无限制的提供新资源(line 66)
	close(p.resources)

	// 清空通道里的资源, 阻止有 goroutine 继续从 p.resources 获取资源进行无
	// 的操作(line 50)
	// 如果需要等待使用完已经存入通道里的资源, 则不需要清空, 只需要关闭通道
	for r := range p.resources {
		r.Close()
	}
}

// Release 和 Close 方法中的互斥量是同一个互斥量，这样可以
// 阻止这两个方法在不同的 goroutine 里同时运行
// 1. 对 closed 的访问存在数据竞态
// 2. 应该保证不往一个已经关闭的通道里发送数据

package main

// 展示使用 pool 来共享一组模拟的
// 数据库连接

import (
	"io"
	"log"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

const (
	maxGoroutines = 25
	poolResources = 2 // 池中的资源数
)

// dbConnection 模拟要共享的资源
type dbConnection struct {
	ID int32
}

// Close 实现了 io.Closer 接口， 以便
// dbConnection 可以被池管理， Close
// 用来完成任意资源的释放管理
func (dbconn dbConnection) Close() error {
	log.Println("Close: Connection", dbconn.ID)
	return nil
}

// 为每个连接分配唯一的id
var idCounter int32

// 工厂函数，当调用一个新连接时，
// 资源池会调用这个函数
func createConnection() (io.Closer, error) {
	id := atomic.AddInt32(&idCounter, 1) // 避免数据竞态
	log.Println("Create: New Connection", id)

	return &dbConnection{id}, nil
}

func main() {
	var wg sync.WaitGroup // 计数信号量
	wg.Add(maxGoroutines)

	// 创建用来管理连接的池
	p, err := New(createConnection, poolResources)
	if err != nil {
		log.Println(err)
	}

	// 使用池中的连接完成查询
	for query := 0; query < maxGoroutines; query++ {
		// 用 maxGoroutines 个 goroutine 并行使用 poolResources 个资源
		go func(q int) { // 匿名函数
			performQueries(q, p) // 对 p.resources 的接收和发送
			wg.Done()
		}(query)
	}
	// 等待 goroutine 结束
	wg.Wait()

	// 关闭池
	log.Println("Shutdown Program.")
	p.Close() // 对 p.resources 关闭且清空
}

func performQueries(query int, p *Pool) {
	// 从池里请求一个连接
	conn, err := p.Acquire()
	if err != nil {
		log.Println(err)
		return
	}

	// 将连接释放回池里
	defer p.Release(conn)

	// 用等待模拟查询响应
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	log.Printf("QID[%d] CID[%d]\n", query, conn.(*dbConnection).ID) // todo
}

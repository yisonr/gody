package main

// todo: 实验效果不能明显观察到并发的效果 ???

// memo3.go 的改进版，在理想情况下，应避免
// 两个或多个 goroutine 几乎同时去获取同一个url

// 上述方案(功能)称为重复抑制(duplicate  suppression)
// 此版本中, map的每个元素是一个指向entry结构的指针，
// 除了包含包含缓存的函数调用结果外，每个entry新增了
// 一个通道 ready, 在设置entry的result字段后，通道会关闭，
// 正在等待的 goroutine 会收到广播, 然后就可以从entry
// 读取结果

//
// 现在调用 get 会先获取保护 cache map 的互斥锁，
// 再从 map 中查询一个指向已有 entry 的指针，
// 如果没有查找到，就分配一个新的 entry ， 最后
// 释放锁。如果要查询的 entry 存在，那么它的值
// 可能还没准备好(另外一个 goroutine 可能正在调用
// 慢函数 f), 所以主调 goroutine 就需要等待 entry
// 准备好才能读取 entry 中的 result 数据，具体的实
// 现方法就是从 ready 通道读取数据，这个操作会一
// 直阻塞到通道关闭.

import "sync"

// Memo 缓存了调用 Func 的结果
type Memo struct {
	f     Func
	mu    sync.Mutex
	cache map[string]*entry
}

// Func 是用来记忆的函数类型
type Func func(key string) (interface{}, error)

type result struct {
	value interface{}
	err   error
}

type entry struct {
	res   result
	ready chan struct{} // res 准备好后会被关闭
}

func New(f Func) *Memo {
	return &Memo{f: f, cache: make(map[string]*entry)}
}

func (memo *Memo) Get(key string) (interface{}, error) {
	memo.mu.Lock()
	e := memo.cache[key]
	if e == nil {
		// 对 key 的第一次访问，这个  goroutine 负责计算数据和广播数据
		// 已准备完毕的消息
		e = &entry{ready: make(chan struct{})}
		memo.cache[key] = e
		memo.mu.Unlock()

		e.res.value, e.res.err = memo.f(key)
		close(e.ready) // 广播数据已准备完毕的消息
	} else {
		// 对这个 key 的重复访问
		memo.mu.Unlock()
		<-e.ready // 等待数据准备完毕
	}

	return e.res.value, e.res.err
}

// entry 中的变量 e.res.value, e.res.err 被多个 goroutine 共享
// 创建 entry 的 goroutine 设置了这两个变量的值，其他 goroutine
// 在收到数据准备完毕的广播后开始读这两个变量，尽管变量被多个
// goroutine 访问， 但此处不需要加上互斥锁。
// 因为: ready 通道的关闭先于其他 goroutine 收到广播事件，所以第一
// 个 goroutine 的变量写入事件也先于后续多个 goroutine 的读取事件
// 在此情况下， 数据竞态不存在(类似于读写锁)

// === RUN   TestConcurrent
// https://music.163.com, 461.534336ms, 135177 bytes
// https://music.163.com, 461.348833ms, 135177 bytes
// https://baidu.com, 504.141624ms, 301477 bytes
// https://baidu.com, 503.584504ms, 301477 bytes
// https://baidu.com, 503.784808ms, 301477 bytes
// https://baidu.com, 503.335049ms, 301477 bytes
// https://bing.com, 904.836655ms, 115494 bytes
// https://bing.com, 904.330417ms, 115494 bytes
// https://bing.com, 904.11781ms, 115494 bytes
// https://bing.com, 903.997998ms, 115494 bytes
// --- PASS: TestConcurrent (0.91s)
// PASS
// ok  	command-line-arguments	0.926s

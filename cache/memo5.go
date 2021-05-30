package main

// memo4.go 中使用一个互斥量来保护被多个调用Get的
// goroutine 访问的 map 变量
// 此版本中新增一种新的设计:
// map 变量限制在一个监控 goroutine 中，而Get调用者
// 则不得不改为发送消息.

//
//
// 使用两种方案构建并发结构: 共享变量并上锁，或者通信
// 顺序进程(communciating sequential process), 两者也挺复杂
// 在给定的情况下很难判断两种方案的优劣，但可以了解这两种
// 方案的对照关系.

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

// request是一条通道消息，key需要用func来调用
type request struct {
	key      string
	response chan<- result // 客户端需要单个 result
}

type Memo struct{ requests chan request }

// New 返回f的函数记忆，客户端之后需要调用close
func New(f Func) *Memo {
	memo := &Memo{requests: make(chan request)}
	go memo.server(f)
	return memo
}

func (memo *Memo) Get(key string) (interface{}, error) {
	response := make(chan result)
	memo.requests <- request{key, response} // 把通道发过去，等待接收数据
	res := <-response
	return res.value, res.err
}

func (memo *Memo) Close() { close(memo.requests) }

func (memo *Memo) server(f Func) {
	cache := make(map[string]*entry)
	for req := range memo.requests {
		e := cache[req.key]
		if e == nil {
			// 对key的第一次请求
			e = &entry{ready: make(chan struct{})}
			cache[req.key] = e
			go e.call(f, req.key) // 调用f(key)
		}
		go e.deliver(req.response)
	}
}

func (e *entry) call(f Func, key string) {
	// 执行函数
	e.res.value, e.res.err = f(key)
	// 通知数据已准备完毕
	close(e.ready)
}

func (e *entry) deliver(response chan<- result) {
	// 等待数据准备完毕
	<-e.ready
	// 向客户端发送结束
	response <- e.res
}

// call() 和 deliver() 方法都需要在独立的 goroutine
// 中运行，以确保监控 goroutine  能持续处理新需求

// go test -run=TestConcurrent -race -v main_test.go requirement.go memo5.go
// === RUN   TestConcurrent
// https://music.163.com, 493.877099ms, 135177 bytes
// https://music.163.com, 494.494466ms, 135177 bytes
// https://baidu.com, 530.553873ms, 301877 bytes
// https://baidu.com, 529.64418ms, 301877 bytes
// https://baidu.com, 530.426674ms, 301877 bytes
// https://baidu.com, 530.246902ms, 301877 bytes
// https://bing.com, 760.111957ms, 115508 bytes
// https://bing.com, 760.347199ms, 115508 bytes
// https://bing.com, 759.412775ms, 115508 bytes
// https://bing.com, 759.800776ms, 115508 bytes
// --- PASS: TestConcurrent (0.76s)
// PASS
// ok  	command-line-arguments	0.784s

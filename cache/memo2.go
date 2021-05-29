package main

import "sync"

// Memo 缓存了调用 Func 的结果
type Memo struct {
	f     Func
	mu    sync.Mutex
	cache map[string]result
}

// Func 是用来记忆的函数类型
type Func func(key string) (interface{}, error)

type result struct {
	value interface{}
	err   error
}

func New(f Func) *Memo {
	return &Memo{f: f, cache: make(map[string]result)}
}

// 并发安全
func (memo *Memo) Get(key string) (interface{}, error) {
	memo.mu.Lock()
	res, ok := memo.cache[key]
	if !ok {
		res.value, res.err = memo.f(key)
		memo.cache[key] = res
	}
	memo.mu.Unlock()
	return res.value, res.err
}

// 虽然进行并行运行测试，不存在数据竞态，
// 但是由于每次调用f时都上锁，因此, Get 把
// 本要并行的I/O操作串行化了
// 在memo3中实现出一个非阻塞的缓存，一个不会把
// 它需要记忆的函数串行运行的缓存

// todo: 既然是串行的为什么没有使用缓存???
// === RUN   TestConcurrent
// https://baidu.com, 514.130776ms, 301624 bytes
// https://music.163.com, 804.729119ms, 135177 bytes
// https://bing.com, 1.448518044s, 115472 bytes
// https://baidu.com, 1.44848424s, 301624 bytes
// https://bing.com, 1.448496793s, 115472 bytes
// https://baidu.com, 1.448309772s, 301624 bytes
// https://music.163.com, 1.448318456s, 135177 bytes
// https://bing.com, 1.448262335s, 115472 bytes
// https://baidu.com, 1.448239783s, 301624 bytes
// https://bing.com, 1.448200147s, 115472 bytes
// --- PASS: TestConcurrent (1.45s)
// PASS
// ok  	command-line-arguments	1.470s

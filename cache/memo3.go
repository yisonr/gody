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

// 并发安全, 且提高了性能, 避免串行调用
// 主调 goroutine 分两次获取锁：第一次用于查询，第二次用于
// 在查询无返回结果时进行更新。在两次加锁的期间(过渡期),
// 其他的 goroutine 也可以使用缓存
func (memo *Memo) Get(key string) (interface{}, error) {
	memo.mu.Lock()
	res, ok := memo.cache[key]
	memo.mu.Unlock()
	if !ok {
		res.value, res.err = memo.f(key)

		// 在两个临界区域之前，可能会有多个 goroutine 来计算 f(key)
		// 并且更新map
		memo.mu.Lock()
		memo.cache[key] = res
		memo.mu.Unlock()
	}
	return res.value, res.err
}

// page 216
//
// 性能得到提升，但仍旧似乎没有使用到缓存????
// 在两个或多个 goroutine 几乎同时调用 Get 来获取
// 同一个url时会出现某些url被获取了多次,
// 两个 goroutine 都首先查询，发现缓存中没有需要的
// 数据， 然后调用那个慢函数f, 最后又都用获得的结果
// 来更新map,其中一个结果会被另外一个覆盖.????
// incomingURLs 函数保证的获取url的同步性啊，为什么
// 会有两个不同的 goroutine 同时查询同一个url???

// === RUN   TestConcurrent
// https://music.163.com, 469.401179ms, 135177 bytes
// https://music.163.com, 478.821934ms, 135177 bytes
// https://baidu.com, 512.351152ms, 301679 bytes
// https://baidu.com, 517.868024ms, 301647 bytes
// https://baidu.com, 521.021186ms, 301645 bytes
// https://baidu.com, 520.114088ms, 301665 bytes
// https://bing.com, 779.193682ms, 115472 bytes
// https://bing.com, 793.557204ms, 115472 bytes
// https://bing.com, 807.563949ms, 115472 bytes
// https://bing.com, 837.135624ms, 115472 bytes
// --- PASS: TestConcurrent (0.84s)
// PASS
// ok  	command-line-arguments	0.866s

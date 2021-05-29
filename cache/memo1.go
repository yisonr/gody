package main

// 函数记忆问题，即缓存函数的结果，达到多次调用但只需
// 计算一次的效果， 解决方案是并发安全的，且会避免
// 简单地对整个缓存使用单个锁带来的锁争夺问题。

// Memo 缓存了调用 Func 的结果
type Memo struct {
	f     Func
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

// 非并发安全
func (memo *Memo) Get(key string) (interface{}, error) {
	res, ok := memo.cache[key]
	if !ok {
		res.value, res.err = memo.f(key)
		memo.cache[key] = res
	}
	return res.value, res.err
}

// 引起了数据竞争
// === RUN   TestConcurrent
// https://music.163.com, 457.234484ms, 135177 bytes
// ==================
// WARNING: DATA RACE
// Write at 0x00c000116d50 by goroutine 22:
//   runtime.mapassign_faststr()
//       /usr/local/go/src/runtime/map_faststr.go:202 +0x0
//   command-line-arguments.(*Memo).Get()
//       /home/lucas/gogo/gody/cache/memo1.go:30 +0x1cd
//   command-line-arguments.TestConcurrent.func1()
//       /home/lucas/gogo/gody/cache/main_test.go:31 +0xdc

// Previous write at 0x00c000116d50 by goroutine 11:
//   runtime.mapassign_faststr()
//       /usr/local/go/src/runtime/map_faststr.go:202 +0x0
//   command-line-arguments.(*Memo).Get()
//       /home/lucas/gogo/gody/cache/memo1.go:30 +0x1cd
//   command-line-arguments.TestConcurrent.func1()
//       /home/lucas/gogo/gody/cache/main_test.go:31 +0xdc

// Goroutine 22 (running) created at:
//   command-line-arguments.TestConcurrent()
//       /home/lucas/gogo/gody/cache/main_test.go:28 +0x1a4
//   testing.tRunner()
//       /usr/local/go/src/testing/testing.go:1123 +0x202

// Goroutine 11 (finished) created at:
//   command-line-arguments.TestConcurrent()
//       /home/lucas/gogo/gody/cache/main_test.go:28 +0x1a4
//   testing.tRunner()
//       /usr/local/go/src/testing/testing.go:1123 +0x202
// ==================
// ==================
// WARNING: DATA RACE
// Write at 0x00c0002aa088 by goroutine 22:
//   command-line-arguments.(*Memo).Get()
//       /home/lucas/gogo/gody/cache/memo1.go:30 +0x1eb
//   command-line-arguments.TestConcurrent.func1()
//       /home/lucas/gogo/gody/cache/main_test.go:31 +0xdc

// Previous write at 0x00c0002aa088 by goroutine 11:
//   command-line-arguments.(*Memo).Get()
//       /home/lucas/gogo/gody/cache/memo1.go:30 +0x1eb
//   command-line-arguments.TestConcurrent.func1()
//       /home/lucas/gogo/gody/cache/main_test.go:31 +0xdc

// Goroutine 22 (running) created at:
//   command-line-arguments.TestConcurrent()
//       /home/lucas/gogo/gody/cache/main_test.go:28 +0x1a4
//   testing.tRunner()
//       /usr/local/go/src/testing/testing.go:1123 +0x202

// Goroutine 11 (finished) created at:
//   command-line-arguments.TestConcurrent()
//       /home/lucas/gogo/gody/cache/main_test.go:28 +0x1a4
//   testing.tRunner()
//       /usr/local/go/src/testing/testing.go:1123 +0x202
// ==================
// https://music.163.com, 474.957942ms, 135177 bytes
// https://baidu.com, 529.91813ms, 301680 bytes
// https://baidu.com, 529.83703ms, 301660 bytes
// https://baidu.com, 533.333429ms, 301642 bytes
// https://baidu.com, 536.285575ms, 301624 bytes
// https://bing.com, 774.299401ms, 115522 bytes
// https://bing.com, 783.739077ms, 115472 bytes
// https://bing.com, 802.320395ms, 115494 bytes
// https://bing.com, 854.660934ms, 115500 bytes
//     testing.go:1038: race detected during execution of test
// --- FAIL: TestConcurrent (0.86s)
// === CONT
//     testing.go:1038: race detected during execution of test
// FAIL
// FAIL	command-line-arguments	0.866s
// FAIL

package main

import "sync"

// RLock 仅可用于在临界区域内对共享变量无写操作的
// 情形, 不应该假设逻辑上只读的函数和方法不会更新
// 一些变量。 比如:
// 一个看起来只是简单访问器的方法可能会递增内部使用
// 的计数器，或者更新一个缓存来让重复的调用更快。

// 仅在绝大部分 goroutine 都在获取锁并且
// 锁竞争比较激烈时(即 goroutine 一般都需要
// 等待后才能获到锁,  RWMutex 才有优势, 因为
// RWMutex 需要更复杂的内部簿记工作，所以在
// 竞争不激烈时比普通互斥锁慢

var mu sync.RWMutex //  读写互斥锁
var balance int

func Balance() int {
	mu.RLock() // 读锁(多读单写锁)
	defer mu.RUnlock()
	return balance
}

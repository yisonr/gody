package main

import (
	"sync/atomic"
	"time"
)

// 原子操作配合互斥锁可以实现高效的单件模式，
// 互斥锁的代价比普通整数的原子读写高很多,
// 在性能敏感的地方可以增加一个数字型的标志位,
// 通过原子检测标志位状态降低互斥锁的使用次数来提高性能

// go 语言高级编程, Page42
// todo: 单件模式
// todo: sync.Once 的实现

// sync/atomic 包对基本数值类型及复杂对象的读写都提供了
// 原子操作的支持， atomic.Value 原子对象提供了 Load() 和 Store() 两个
// 原子方法，分别用于加载和保存数据，返回值和参数都是 interface{} 类型,
// 因此用于任意的自定义复杂类型。
//

var config atomic.Value

func main() {
	// 加载配置
	config.Store(loadConfig())

	// 启动一个后台线程, 加载更新后的
	// 配置信息(一个goroutine分配到一个逻辑处理器上执行)
	go func() {
		for {
			time.Sleep(time.Second) // 1秒更新一次配置
			config.Store(loadConfig())
		}
	}()

	for i := 0; i <= 10; i++ {
		go func() {
			for r := range requests() {
				c := config.Load()
				// ...
			}
		}()
	}
}

func loadConfig() {
	// 加载配置
	// ...
}

func requests() <-chan int {
	requs := make(chan int, 4)
	return requs
}

/*
* 以上是一个简化的生产者消费者模型:
* 后台线程生成最新的配置信息，前台多个工作者线程线程获取最新的配置信息，
* 所有线程共享配置信息资源.
*
 */

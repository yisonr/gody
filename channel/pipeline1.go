package main

import (
	"fmt"
	"time"
)

func main() {
	naturals := make(chan int)
	squarers := make(chan int)

	// counter
	go func() {
		for x := 0; x < 10; x++ {
			naturals <- x // 发送
		}
		close(naturals)
		// 结束时，必需关闭每一个通道
		// 只有通知接收方 goroutine 所有的数据都发送完毕的
		// 时候才需要关闭通道，通道也是可以通过垃圾回收器
		// 根据它是否可访问来决定是否回收它，而不是根据它是否关闭
	}()

	// squarer
	go func() {
		for x := range naturals {
			// x := <-naturals   // 接收, 在发送前阻塞
			squarers <- x * x // 发送
		}
		close(squarers) // 关闭后发送导致宕机,
		// 接收所有已经发送的数据, 取完后得到元素类型的零值
	}()

	// printer(在主 goroutine 中)
	for x := range squarers {
		fmt.Println(x) // 接收，在发送前阻塞
		time.Sleep(1000 * time.Millisecond)
	}
}

package main

//TODO: time.Tick 函数返回的只接收
// 通道的容量是多少???

import (
	"fmt"
	"os"
	"time"
)

// select 使用1s等待中止，但不会更长

func main() {
	abort := make(chan struct{})
	go func() {
		os.Stdin.Read(make([]byte, 1))
		abort <- struct{}{}
	}()

	fmt.Println("Commencing countdown. Press return to abort.")
	tick := time.Tick(1 * time.Second) // 这个通道的容量是多少???
	// time.Tick 的行为很像创建一个goroutine在循环中
	// 调用time.Sleep，然后在它每次醒来时发送事件.
	// 当倒计时函数返回时，它停止从tick通道中接收事件,
	// 但是计时器goroutine还在运行, 会不断的向没有goroutine
	// 接收的通道发送，这就引起了goroutine 泄露.
	for countdown := 10; countdown > 0; countdown-- {
		fmt.Println(countdown)
		select {
		case <-tick:
			// 不执行操作
		case <-abort:
			fmt.Println("Launch aborted!")
			return
		}
	}
	fmt.Println("Start launch!")
	// launch()
}

/*
 TODO
 time.Tick 函数很方便使用，但是它仅仅在应用的整个
 生命周期中都需要才合适。 否则,应使用以下模式:
 ticker := time.NewTicker(1*time.Second)
 <-ticker.C  // 从 ticker 的通道接收
 ticker.Stop() // 造成 ticker 的goroutine终止
*/

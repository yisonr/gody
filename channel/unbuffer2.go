package main

/*
无缓冲通道
模拟4个 goroutine 间的接力比赛

在每一个跑者 goroutine 中，首先等待拿到接力棒
(第一个跑者的接力棒由主 goroutine 传输), 创建
下一位跑者, 并开启下一位跑者 goroutine 等待接力，
开始绕跑道跑一段时间(time.Sleep)
每一个跑者 goroutine 中都应判断自己是不是最后一
个接力跑者，若是则比赛结束， 程序退出，若不是，
则把接力棒通过无缓冲通道传给等待接力的跑者goroutine

无缓冲通道保证了接力棒在4位跑者 goroutine 中的
同步传递
*/

import (
	"fmt"
	"sync"
	"time"
)

// 等待程序结束
var wg sync.WaitGroup

func main() {
	baton := make(chan int) // 接力棒
	// 为最后以为跑步者将计数加1
	wg.Add(1)
	// 因为接力赛中的跑者使用接力棒的操作是
	// 同步操作， 所有只需要等待一个 goroutine 的结束
	// 即最后一位跑者 goroutine 的结束

	// 第一位跑步者持有接力棒
	go Runner(baton)

	// 开始比赛
	baton <- 1

	// main 函数阻塞在此，等待比赛结束
	wg.Wait()
}

// Runner 模拟接力比赛中的一位跑步者
func Runner(baton chan int) {
	var newRunner int

	// 等待接力棒
	runner := <-baton // 1号选手通过主 goroutine 来到第一个

	// 开始绕着跑道跑步
	fmt.Printf("Runner %d running with baton\n", runner)

	// 创建下一位跑步者
	if runner != 4 {
		newRunner = runner + 1
		fmt.Printf("Runner %d to the line\n", runner)
		go Runner(baton)
	}

	// 围绕跑道跑
	time.Sleep(100 * time.Millisecond)

	// 比赛结束了?
	if runner == 4 {
		fmt.Printf("Runner %d finished, race over\n", runner)
		// 第4个跑步者完成了比赛
		wg.Done()
		return
	}

	// 将接力棒交给下一位跑者
	fmt.Printf("Runner %d exchange with runner %d\n", runner, newRunner)
	baton <- newRunner // 锁住，直到交接完成
}

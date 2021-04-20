package main

/*
 无缓冲通道
 模拟网球比赛中，两位选手会把球在两人之间来回传递
 选手总是存在一下两种状态:
 1.等待接球
 2.将球打向对方
用无缓冲通道模拟球的来回
*/

// 主 goroutine 只进行通道发送，第一次发球开始游戏
// 然后球在 goroutine 之间流动，接住球的 goroutine
// 继续传出球，没接住球的 goroutine 输掉比赛

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var wg sync.WaitGroup

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	court := make(chan int)
	wg.Add(2) // 等待 goroutine 完成

	// 启动两个选手
	go player("Nadal", court)
	go player("Djokovic", court)

	// 发球
	court <- 1

	// 等待游戏结束
	wg.Wait()
}

// 模拟一个选手在打网球
func player(name string, court chan int) {
	// 在函数退出时调用 Done 来通知 main 函数工作已经完成
	defer wg.Done()

	for {
		// 等待球被击打过来
		ball, ok := <-court
		if !ok {
			// 如果通道关闭，赢
			fmt.Printf("Player %s won\n", name)
			return
		}

		// 选随机数，然后用这个数来判断是否丢球
		n := rand.Intn(100)
		if n%3 == 0 {
			fmt.Printf("Player %s Missed\n", name)

			// 关闭通道，输了
			close(court)
			return
		}

		// 显示击球数，并将击球数加1
		fmt.Printf("Player %s hit %d\n", name, ball)
		ball++

		// 将球打向对手
		court <- ball
	}
}

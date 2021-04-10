package main

import (
	"fmt"
	"time"
)

// 通道ch的缓冲区大小为1, 要么是空的，
// 要么是满的，因此只有其中一个情况下可执行

func main() {
	ch := make(chan int, 1)
	// 增加缓冲区的容量，会使得输出变得不可确定
	// 因为当缓冲既不空也不满的情况，相当于select
	// 语句在扔硬币做选择
	for i := 0; i < 100; i++ {
		time.Sleep(1 * time.Second)
		select {
		// 如果多个情况同时满足，select会随机选择一个,
		// 保证每一个通道有相同的机会被选中
		case x := <-ch:
			fmt.Println(x)
		case ch <- i:
			//... 不执行操作
		}
	}
}

/*
	有时候需要在一个通道上发送或接收，但是不想在
	通道没有准备好的情况下被阻塞--非阻塞通信。这使用
	select 也可以做到。
	select 可以有一个默认情况，它用来指定在没有其他
	的通信发生可以立即执行的动作.

	下面的语句尝试从 abort 通道中接收一个值，如果没有
	值，什么也不做。这是一个非阻塞的接收操作，重复这个
	动作称为对通道轮询。
	------ 重复才称为轮询
	select{
	case<-abort:
		fmt.Printf("Launch aborted!\n")
		return
	default:
		// 不执行任何操作
	}
	------
*/

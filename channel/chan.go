package main

import "fmt"

func main() {
	// ch := make(chan int) // 无缓冲通道，发送操作将会阻塞, 造成死锁(deadlock)
	ch := make(chan int, 2)
	ch <- 4 // 发送给通道
	ch <- 5
	x := 0
	x = <-ch // 从通道接收
	fmt.Println(x)
	close(ch)
	x = <-ch
	x = <-ch
	x = <-ch
	// 关闭 close 设置一个标志位来指示值当前以及发送(给通道)完毕,
	// 这个通道后面无值: 关闭后的发送操作导致宕机。
	// 关闭后的通道上进行接收操作(从通道接收), 将获取已经发送的值，
	// 直到通道为空， 这时任何接收操作会立即完成，同时获取到一个通道
	// 元素类型对应的零值。
	fmt.Println(x)
}

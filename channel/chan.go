package main

import "fmt"

func Demo() {
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

func main() {
	ch := make(chan string, 3)
	// 缓冲通道队列
	// 先进先出, 队列满时，发送操作阻塞
	// 队列空时，接收操作阻塞
	ch <- "string" // 发送
	ch <- "test"
	ch <- "A people"
	close(ch)
	<-ch                 // 接收
	fmt.Println(cap(ch)) // 内置的cap获取通道缓冲区的容量
	fmt.Println(len(ch)) // 内置的len获取当前通道内的元素个数
	// 但不能粗暴的将缓冲通道作为队列在单个 goroutine 中使用，
	// 通道和 goroutine 的调度深度关联，如果，没有另外一个
	//  goroutine 从通道进行接收，发送者(也许是整个程序) 有被永久
	// 阻塞的风险。

}

package main

import "fmt"

// 当一个通道用做函数的形参时，它会被有意的限制不能发送或不能接收
// 将这种意图文档化可以避免误用，go的类型系统提供了“单向通道类型”,
// 仅仅导出发送或接收操作
// chan<- int 只能发送的通道(int类型)
// <-chan int 只能接收的通道(int类型)
// close 操作说明了通道上没有数据在发送，
// 仅仅在发送方 goroutine 上才能调用，所以
// 试图关闭一个仅能接收的通道会在编译时报错

func counter(out chan<- int) {
	for x := 0; x < 10; x++ {
		out <- x // 发送
	}
	close(out)
}

func squarer(out chan<- int, in <-chan int) {
	for x := range in {
		// x := <-naturals   // 接收, 在发送前阻塞

		out <- x * x // 发送
	}
	close(out) // 关闭后发送导致宕机,
	// 接收所有已经发送的数据, 取完后得到元素类型的零值
}

// func printer(in chan int) { // 也可以传入双向通道
func printer(in <-chan int) {
	for x := range in {
		fmt.Println(x) // 接收，在发送前阻塞
	}
}

func main() {
	naturals := make(chan int)
	squarers := make(chan int)
	go counter(naturals)
	go squarer(squarers, naturals)
	printer(squarers)
}

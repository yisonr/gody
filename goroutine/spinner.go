package main

import (
	"fmt"
	"time"
)

func main() {
	go spinner(100 * time.Millisecond)
	// main 返回或程序退出时，此 goroutine 直接被暴力终结
	// 除此之外没有程序化的方法让一个 goroutine 来停止另一个
	const n = 45
	fibN := fib(n) // slow
	fmt.Printf("\rFibonacci(%d) = %d\n", n, fibN)
}

func spinner(delay time.Duration) {
	for {
		for _, r := range `-\|/` {
			fmt.Printf("\r%c%c", r, r)
			time.Sleep(delay)
		}
	}
}

// 此递归调用效率极低，因为其进行了重复运算
func fib(x int) int {
	if x < 2 {
		return x
	}
	return fib(x-1) + fib(x-2)
}

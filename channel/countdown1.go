package main

// 火箭发射程序

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Commencing countdown.")
	tick := time.Tick(1 * time.Second) // 只接收通道
	for countdown := 10; countdown > 0; countdown-- {
		fmt.Println(countdown)
		<-tick
	}
	// launch()
}

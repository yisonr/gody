package main

import "fmt"

// 顺序一致性内存模型

var a string
var done bool

func setup() {
	a = "hello, world"
	done = true
}

func count() {
	go setup()
	for !done {
	}
	fmt.Println(a)
}

func main() {
	for i := 1; i <= 1000; i++ {
		count()
	}
}

package main

// 编译器并不是总能自动获得一个值的地址
// 方法集: todo

import "fmt"

type duration int

func (d *duration) pretty() string {
	return fmt.Sprintf("Duration: %d", *d)
}

func main() {
	fmt.Println(duration(42).pretty())
	//compiler: pretty in not in method set  of duration
}

// ./addr.go:13:26: cannot call pointer method on duration(42)
// ./addr.go:13:26: cannot take the address of duration(42)

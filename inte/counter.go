package main

// interface 实际上是一个约定
// 函数中的的形参(interface) 约定实参要实现的方法
// 实现该 interface 的类型作为实参
// 在函数中具体化对实参方法的调用

import (
	"bufio"
	"fmt"
	"strings"
)

type ByteCounter int

func (c *ByteCounter) Write(p []byte) (int, error) {
	*c += ByteCounter(len(p)) // 类型转换
	return len(p), nil
}

type WordCounter int

func (c *WordCounter) Write(p []byte) (int, error) {
	scanner := bufio.NewScanner(strings.NewReader(string(p)))
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		*c += 1
	}
	return len(p), nil
}

type LineCounter int

func (c *LineCounter) Write(p []byte) (int, error) {
	scanner := bufio.NewScanner(strings.NewReader(string(p)))
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		*c += 1
	}
	return len(p), nil
}

func main() {
	// var c ByteCounter
	// c.Write([]byte("hello"))
	// fmt.Println(c)

	// c = 0 // 重置计数器
	// var name string = "sdfsf"
	// fmt.Fprintf(&c, "hello, %s", name)
	// fmt.Println(c)
	var l LineCounter
	l.Write([]byte(`sdfsdfsdfsf
				sdfnsdfds
				benchmark
				lesi
				license
				sdsfs`))
	fmt.Println(l)

	l = 0 // 重置计数器
	var name string = `sdfsf
					   word
					   lalala`
	fmt.Fprintf(&l, "hello, %s", name)
	fmt.Println(l)
}

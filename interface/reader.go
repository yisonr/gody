package main

import (
	// "fmt"
	"io"
	"os"
)

func LimitReader(r io.Reader, n int64) io.Reader {
	// Reader 抽象了所有可以读取字节的类型
	// 此函数在读取字节的操作上进行控制
	// 一个具体类型实现 Reader 抽象类型的接口
	r.Read([]byte("sdfsfsf"))
	return
}

func main() {
	LimitReader(os.Stdin, 10)
}

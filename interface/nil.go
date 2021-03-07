package main

// 空接口值和仅仅动态值为 nil 的接口值的区别
// 动态分发机制
// 对于某些类型(*os.File)空接收值是合法的
// 但对于 *bytes.Buffer 的空接收值在尝试
// 访问缓冲区的时候会崩溃

import (
	"bytes"
	"fmt"
	"io"
)

const debug = false

func main() {
	var buf *bytes.Buffer
	if debug {
		buf = new(bytes.Buffer) // 启用输出收集
	}
	f(buf)
	fmt.Println(buf)
}

func f(out io.Writer) {
	if out != nil {
		out.Write([]byte("done"))
	}
}

package main

// 空接口类型 interface{}
// 对其实现类型没有任何要求,
// 可以把任何值赋给空接口类型

// 依靠空接口类型, fmt.Println, errorf 函数
// 才能接受任意类型的参数

import (
	"bytes"
	"fmt"
	"io"
)

func main() {
	var any interface{}
	any = true
	any = 12.34
	any = new(bytes.Buffer)

	// 如下声明在编译器就断言了 *byte.Buffer 类型
	// 的一个值必然实现了 io.Writer
	// *bytes.Buffer 必须实现 io.Writer
	var w io.Writer = new(bytes.Buffer)
	// *bytes.Buffer 的任意值都实现了io.Writer 接口
	// 甚至 nil, 用(*bytes.Buffer)(nil) 做强制类型转换后，
	// 也实现了这个接口
	fmt.Println(w)
}

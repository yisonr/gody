// 打印命令行参数
// - 如果变量在声明的时候没有明确的初始化，
// - 将隐式地初始化为这个类型的空值
// - := 用于短变量声明
// - i++ 是语句不是表达式，且仅支持后缀(j = i++, ++i 都是不合法的)
// - _(下划线)空标识符表示临时变量
package main

import (
	"fmt"
	"os"
)

func main() {
	var s, sep string
	for i := 1; i < len(os.Args); i++ {
		s += sep + os.Args[i]
		sep = " "
	}
	fmt.Println(s)
	// for index, value := range os.Args {
	// 	fmt.Printf("%d, %s\n", index, value)
	// }
}

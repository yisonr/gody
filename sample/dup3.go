// 找出文件或标准输入的重复行
// 一次读取整个输入到大块内存，一次性地分割所有行
// for range 可以遍历数组、切片、字符串、map 及通道（channel）
// 要注意的是，val 始终为集合中对应索引的值拷贝,
// 因此它一般只具有只读性质，对它所做的任何修改都不会影响到集合中原有的值
// ReadFile 读取整个文件的内容
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	counts := make(map[string]int)
	for _, filename := range os.Args[1:] {
		data, err := ioutil.ReadFile(filename)
		// data 是一个可以转化成字符串的字节slice
		if err != nil {
			fmt.Fprintf(os.Stderr, "dup3: %v\n", err)
			continue
		}
		for _, line := range strings.Split(string(data), "\n") {
			counts[line]++
		}
	}
	for line, n := range counts {
		if n > 1 {
			fmt.Printf("%d\t%s\n", n, line)
		}
	}

}

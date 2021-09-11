package main

import (
	"fmt"
	"unsafe"
)

/*
 * bool 为什么占用一个字节, 而不是一位
 * 因为内存的最小单元是字节
 *
 *
 */

func main() {
	// 声明bool变量
	b := true
	fmt.Printf("%d", unsafe.Sizeof(b))
	// 强转内存对齐
	fmt.Printf("bool value = %d\n", *(*uint8)(unsafe.Pointer(&b)))

}

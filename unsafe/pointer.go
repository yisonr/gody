package main

import (
	"fmt"
	"unsafe"
)

/*
	unsafe.Pointer 类型是一种特殊类型的指针, 它可以存储任何变量的地址;
	unsafe.Pointer 类型的指针可以比较并且可以和 nil 做比较;

	一个普通的指针 *T 可以转换为 unsafe.Pointer 类型的指针, 一个 unsafe.Pointer
	类型也可以转换回普通指针, 而且不必和原来的类型 *T 相同;

	unsafe.Pointer 也可转换为 uintptr 类型, uintptr 类型保证了指针所指向地址的
	数值, 便可对地址进行数值计算(uintptr 类型是一个足够大的无符号整形, 可以
	用来表示任何地址); 从 uintptr 到 unsafe.Pointer 的转换也会破坏类型系统,
	因为并不是所有的数值都是合法的内存地址;

	很多 unsafe.Pointer 类型的值都是从普通指针到原始内存地址以及再从内存地址
	到普通指针进行转换的中间;

*/

// 通过转换一个 *float64 类型的指针到 *uint64 类型, 可以查看浮点类型变量的
// 位模式(TODO:)
func Float64bits(f float64) uint64 {
	// 内存地址未变, 只是给地址赋上了新的类型
	return *(*uint64)(unsafe.Pointer(&f))
}

// 也可以通过结果指针来更新位模式, 这对一个浮点类型的变量来说是无害的, 但是
// 通常使用 unsafe.Pointer 进行类型转换可以将任意值写入内存中, 并因此破坏了
// 类型系统;
// func main() {
// 	fmt.Printf("%#016x\n", Float64bits(1.0)) // 0x3ff0000000000000
// }

func updateFieldValue() {
	var x struct {
		a bool
		b int16
		c []int
	}

	// 获取变量x的地址, 再加上其成员b的地址偏移量, 并将结果转为 *int16 类型,
	// 接着通过这个指针更新 x.b 的值
	pb := (*int16)(unsafe.Pointer(
		uintptr(unsafe.Pointer(&x)) + unsafe.Offsetof(x.b)))
	*pb = 42
	fmt.Println(x.b) // 42
}

// 注意: 不要引入 uintptr 类型的临时变量来破坏整行代码
func errExample() {
	var x struct {
		a bool
		b int16
		c []int
	}
	tmp := uintptr(unsafe.Pointer(&x)) + unsafe.Offsetof(x.b)
	pb := (*int16)(unsafe.Pointer(tmp)) // 编译器警告: possible misuse of unsafe.Pointer
	*pb = 42
	fmt.Println(x.b) // 尽管运行此示例时输出满足预期, 但是结果是未知的
}

/*
	TODO: 垃圾回收器
	垃圾回收器在内存中把变量移来移去以减少内存碎片或者为了进行簿计工作(TODO),
	这种类型的垃圾回收器称为移动的垃圾回收器, 当一个变量在内存中移动后, 该
	变量所指向旧地址的所有指针都需要更新以指向新地址, 从垃圾回收的角度,
	unsafe.Pointer 是一个变量的指针, 当变量移动的时候, 它的值也需要改变,
	而 uintptr 仅仅是一个数值, 所以它的值是不会变的;
	所以如果使用引入 uintptr 类型的临时变量, 会使得垃圾回收器无法通过非指针
	变量 tmp 了解它背后的指针; 当定义 tmp 临时变量后的下一行代码执行时,
	变量 x 可能在内存中已经移动了, 这时 tmp 中的值就不是变量 &x.b 的地址了;
	赋值语句将向任意的内存地址写入值 42;


	goroutine 栈会根据需要增长, 这时, 旧栈上面的所有变量都会重新分配到新的,
	更大的栈上, 所以不能指望变量的地址在它的整个生命周期都不会变;

	当调用一个返回 uintptr 类型的库函数时, 其结果应该立刻转换为 unsafe.Pointer
	来确保它们在接下来的代码中指向同一个变量, 如 reflect 包的函数:
	- func(Value) Pointer() uintptr
	- func(Value) UnsafeAddr() uintptr
	- func(Value) InterfaceData() [2]uintptr
*/

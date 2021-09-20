package main

import (
	"fmt"
	"io"
	"os"
	"reflect"
)

/*
 * 反射功能由 reflect 包提供，两个重要的类型: Type 和 Value
 *
 *  Type: 表示go语言的一个类型，是一个有很多方法的接口，这些方法可以
 *		用来识别类型以及透视类型的组成部分.(如一个结构的各个字段或者
 *      一个函数的各个参数)
 *
 *		reflect.Type 接口只有一个实现，即类型描述符，接口值中的动态类型也
 *		是类型描述符
 *
 *		reflect.TypeOf 函数接受任何的 interface{} 参数，并把接口中的动态类型
 *			以 reflect.Type 形式返回
 *
 *  Value: 可以包含一个任意类型的值
 *
 *		reflect.ValueOf 函数接受任何的 interface{} 参数，并把接口中的动态值
 *			以 reflect.Value 形式返回
 *
 */

func main() {
	// 把一个具体值赋给一个接口类型时会发生一个隐式类型转换，转换会生成
	// 包含两部分内容的接口值:
	// - 动态类型部分是操作数的类型(int)
	// - 动态值部分是操作数的值(3)
	t := reflect.TypeOf(3)
	fmt.Println(t.String()) // int
	fmt.Println(t)          // int

	// 因为 reflect.TypeOf 总是返回一个接口值对应的动态类型, 所以它的返回总是
	// 具体类型(而不是接口类型)
	// 比如以下代码输出的是 "*os.File" 而不是 "io.Writer"
	var w io.Writer = os.Stdout
	fmt.Println(reflect.TypeOf(w)) // *os.File

	// reflect.TypeOf 满足 fmt.Stringer , 因为输出一个接口值的动态类型在调试
	// 和日志中很常用，所以 fmt.Printf 提供了一个简写方式 %T, 内部实现就使用
	// 了 reflect.TypeOf
	fmt.Printf("%T\n", 3) // int

	// 与 reflect.TypeOf 类似，reflect.ValueOf 的返回值也都是具体值，不
	// 过 reflect.Value 也可以包含一个接口值
	v := reflect.ValueOf(3) // 一个 reflect.Value
	fmt.Println(v)          // "3"
	fmt.Printf("%v\n", v)   // "3"
	fmt.Println(v.String()) // 注意: <int Value> , 仅仅暴露了类型
	// reflect.Value 也满足 fmt.Stringer , 但除非 Value 包含的是一个字符串，
	// 否则 String 方法的结果仅仅暴露了类型，通常要使用 fmt 包的 %v 功能，
	// 它对 reflect.Value 会进行特殊处理.
}

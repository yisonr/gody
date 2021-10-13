package main

import (
	"fmt"
	"os"
	"reflect"
)

// 使用 reflect.Value 设置值

// 除过使用 reflect 解析变量值，也可以用来修改值
// 比如 x, x.f[1], *p 这样的表达式表示一个变量，而 x+1，f(2) 之类
// 的表达式则不表示变量，一个变量是一个可寻址的存储区域，其中包含了
// 一个值，并且它的值可以通过这个地址来更新

// 对于 reflect.Value 也有类似的区分，某些是可寻址的，其他的则不可
func main1() {
	x := 2
	fmt.Println(&x)
	a := reflect.ValueOf(2)
	fmt.Println(&a) // <int Value> a 的值不可寻址，它包含的仅仅是整数2的一个副本
	b := reflect.ValueOf(x)
	fmt.Println(&b) // <int Value> b 的值同a
	c := reflect.ValueOf(&x)
	fmt.Println(c)  // 0xc0000140f0  c 的值也不可寻址，包含的是指针 &x 的一个副本
	fmt.Println(&c) // <*int Value>
	//  以上 a, b, c 皆不可寻址

	// 可以通过一个指针间接得到一个可寻址的 reflect.Value,
	// 即使这个指针是不可寻址的
	// d 是通过对 c 中的指针提领得来的， 所有可寻址 ??? 提领???
	// 得到一个可寻址的 reflect.Value()， 如下:
	// - 调用 Addr() ，返回一个Value, 其中包含了一个指向变量的指针
	// - 在 Value 上调用 Interface() 方法，返回一个包含指针的 interface{} 值
	// - 使用类型断言把接口内容转换为一个普通指针
	d := c.Elem()                     // 或 d:= reflect.ValueOf(&x).Elem() d代表变量x
	px := d.Addr().Interface().(*int) // px := &x
	fmt.Println(px)                   // 0xc0000140f0
	fmt.Println(*px)
	*px = 4
	fmt.Println(*px)

	/* TODO: 理解
	 * 可寻址的常见规则都在反射包里有对应项，如:
	 * slice 的脚标表达式e[i]隐式的做了指针去引用，所以即使e是不可寻址的，
	 * 这个表达式仍然是可寻址的，类似，reflect.ValueOf(e).Index(i) 代表
	 * 一个变量，尽管 reflect.ValueOf(e) 是不可寻址的, 这个变量也是可寻址的
	 */

	// 可以通过变量的 CanAddr 方法询问 reflect.Value 变量是否可寻址
	fmt.Println(a.CanAddr()) // "false"
}

func main2() {
	// 还可以通过可寻址的 reflect.Value 更新变量，即直接调用 reflect.Value.Set
	// 方法;
	x := 2
	c := reflect.ValueOf(&x)
	d := c.Elem() // 或 d:= reflect.ValueOf(&x).Elem() d代表变量x
	d.Set(reflect.ValueOf(6))
	fmt.Println(x)
	/*
	 * 平常由编译器来检查的那些可赋值性条件，在这种情况下则是在运行时由 Set
	 * 方法来检查; 上面的变量和值都是 int 类型，但如果变量类型是 int64,
	 * 程序会崩溃，所以需要确保值对于变量类型是可赋值的
	 */
	// d.Set(reflect.ValueOf(int64(5)))
	// panic: reflect.Set: value of type int64 is not assignable to type int

	// 在不可寻址的 reflect.Value 上调用 Set 方法也会崩溃
	// y := 2
	// z := reflect.ValueOf(y)
	// z.Set(reflect.ValueOf(2))
	// panic: reflect: reflect.Value.Set using unaddressable value

	/*
	 * 为一些基本类型特化的 Set 变种: SetInt, SetUint, SetString, SetFloat
	 * 上述方法具有一定的容错性, 只要变量类型是某种带符号的整数，比如 SetInt,
	 * 甚至是底层类型为带符号的命名类型，都可以成功
	 * 如果值太大了会无提示的截断
	 *
	 * 但在指向 interface{} 变量的 reflect.Value 上调用 SetInt 会崩溃, 但可以
	 * 使用 Set
	 */

	var t interface{}
	rt := reflect.ValueOf(&t).Elem()
	// rt.SetInt(2)
	// panic: reflect: call of reflect.Value.SetInt on interface Value
	rt.Set(reflect.ValueOf(3))
	fmt.Println(t)
}

func main() {
	// 通常 go 语言的常规方法无法读取未导出结构字段的值，可以使用反射读取
	// 但反射不能更新这些值
	stdout := reflect.ValueOf(os.Stdout).Elem() // *os.Stdout, 一个 os.File变量
	fmt.Println(stdout.Type())                  // os.File
	fd := stdout.FieldByName("fd")
	// TODO: 和教材上不一样
	// -------------教材内容---------------------- go 程序设计 Page267
	fmt.Println(fd.Int()) // 1
	// fmt.Println(fd.Int()) // panic: reflect: call of reflect.Value.Int on zero Value
	// fd.SetInt(2) // 未导出字段
	// -------------教材内容----------------------
	// fmt.Println(fd.Int()) // panic: reflect: call of reflect.Value.Int on zero Value
	// fd.SetInt(2) // panic: reflect: call of reflect.Value.Int on zero Value

	// TODO: 验证
	// 一个可寻址的 reflect.Value 会记录它是否是通过遍历一个未导出字段来获得的,
	// 如果是未导出字段，则不允许修改
	// 所以在更新变量前，CanAddr 并不能保证正确， CanSet 方法才能正确报告
	// 一个 reflect.Value 是否可寻址且可更改
	fmt.Println(fd.CanAddr(), fd.CanSet()) // true false(教程)
	// panic: reflect: call of reflect.Value.Int on zero Value (实际)
}

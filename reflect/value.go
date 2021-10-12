package main

import (
	"fmt"
	"reflect"
)

// 使用 reflect.Value 设置值

// 除过使用 reflect 解析变量值，也可以用来修改值
// 比如 x, x.f[1], *p 这样的表达式表示一个变量，而 x+1，f(2) 之类
// 的表达式则不表示变量，一个变量是一个可寻址的存储区域，其中包含了
// 一个值，并且它的值可以通过这个地址来更新

// 对于 reflect.Value 也有类似的区分，某些是可寻址的，其他的则不可
func main() {
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

}

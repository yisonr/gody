package main

import (
	"fmt"
	"strconv"
)

/*
 *  有时候需要写一个有能力统一处理各种值类型的函数，而这些类型可能无法共
 *  享同一个接口，也可能布局未知，也可能这个类型在我们设计函数时还不存在,
 *  甚至这个类型会同时存在上面三种问题.
 *
 */

//   模拟实现 fmt.Printf 的格式化逻辑的实现
func Sprint(x interface{}) string {
	type stringer interface {
		String() string
	}
	switch x := x.(type) {
	case stringer:
		return x.String()
	case string:
		return x
	case int:
		return strconv.Itoa(x)
		//... 对 int16, int32 ... 等类型做类似的处理
	case bool:
		if x {
			return "true"
		}
		return "false"
	default:
		// slice, chan, func, map, pointer, struct
		return "???"
	}
}

// 为什么要使用反射:
// 以上使用分支的方法并不能实现 fmt.Printf 函数的功能, 虽然可以根据更多类型的
// 的需要添加更多的分支，但这样的类型有无限种, 更何况还有自己命名的类型, 当
// 无法透视一个未知类型的布局时，这段代码就无法继续，这时就需要反射了
//

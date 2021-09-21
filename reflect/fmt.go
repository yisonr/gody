package main

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
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
// 无法透视一个未知类型的布局时，这段代码就无法继续，这时就需要反射
//

/*
* 基于 ./reflect.go 中对反射的理解，可以不用类型分支，而使用 reflect.Value 的
* Kind 方法区分不同的类型, 尽管有无限种类型，但类型的分类(Kind)只有几种:
*	- 基础类型: Bool, String 以及各种数字类型
*  - 聚合类型: Array, Struct
*  - 引用类型: Chan, Func, Ptr, Slice, Map，接口类型和 interface
*  - Invalid 类型， 表示它们没有任何值(reflect.Value 的零值属于 Invalid 类型)
*
* 下面尝试使用 Kind
 */

// Any  把任何值格式化为一个字符串
func Any(value interface{}) string {
	return formatAtom(reflect.ValueOf(value))
}

// formatAtom 格式化一个值，且不分析它的内部结构
func formatAtom(v reflect.Value) string {
	/*
	 * 该函数把每个值当做一个没有内部结构且不可分割的物体，对于聚合类型(结构体和
	 * 数组) 以及接口, 只输出了值的类型，对引用类型(通道，函数，指针，slice 和
	 * map), 输出了类型和以十六进制表示的引用地址.
	 *
	 */
	switch v.Kind() {
	case reflect.Invalid:
		return "invalid"

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Int64:
		// Int returns v's underlying value, as an int64.
		// FormatInt 函数的第一个参数要求的类型是 int64
		return strconv.FormatInt(v.Int(), 10)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32,
		reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(v.Uint(), 10)

		// 省略浮点数和复数的分支

	case reflect.Bool:
		return strconv.FormatBool(v.Bool())

	case reflect.String:
		return strconv.Quote(v.String())

	case reflect.Chan, reflect.Func, reflect.Ptr, reflect.Slice, reflect.Map:
		return v.Type().String() + " 0x" +
			strconv.FormatUint(uint64(v.Pointer()), 16)

	default: // reflect.Array, reflect.Struct, reflect.Interface
		return v.Type().String() + " value"
	}
}

func main() {
	//  Kind 只关心底层实现， 所以 Any 函数对命名类型的效果也很好
	var x int64 = 1
	var d time.Duration = 1 * time.Nanosecond
	fmt.Println(Any(x))
	fmt.Println(Any(d))
	fmt.Println(Any([]int64{x}))
	fmt.Println(Any([]time.Duration{d}))
}

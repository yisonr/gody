package main

import (
	"fmt"
	"reflect"
)

/*
 * 应尽可能避免在包的API里暴露反射相关的内容
 *
 * 定义一个未导出的函数 display 做真正的递归处理， 再暴露 Display
 * 而 Display 则只是一个简单的封装，并且接收一个 interface{} 参数
 *
 */

// 调试工具函数(组合类型的显示), 给定任意一个复杂值x, 输出复杂值的完整结构,
// 并对找到的每个元素标上相应的路径
func Display(name string, x interface{}) {
	fmt.Printf("Display %s (%T):\n", name, x)
	display(name, reflect.ValueOf(x))
}

func display(path string, v reflect.Value) {
	switch v.Kind() {
	case reflect.Invalid:
		fmt.Printf("%s = invalid\n", path)
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			display(fmt.Sprintf("%s[%d]", path, i), v.Index(i))
		}
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			filePath := fmt.Sprintf("%s.%s", path, v.Type().Field(i).Name)
			display(filePath, v.Field(i))
		}
	case reflect.Map:
		for _, key := range v.MapKeys() {
			display(fmt.Sprintf("%s[%s]", path, formatAtom(key)),
				v.MapIndex(key))
		}
	case reflect.Ptr:
		if v.IsNil() {
			fmt.Printf("%s=nil\n", path)
		} else {
			display(fmt.Sprintf("(*%s)", path), v.Elem())
		}
	case reflect.Interface:
		if v.IsNil() {
			fmt.Printf("%s=nil\n", path)
		} else {
			fmt.Printf("%s.type=%s\n", path, v.Elem().Type())
		}
	default:
		fmt.Printf("%s=%s\n", path, formatAtom(v))
	}
}

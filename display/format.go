package main

import (
	"reflect"
	"strconv"
)

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

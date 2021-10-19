package main

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

// 使用一个反射示例 reflect.Type 显示一个任意值的类型并枚举它的方法

func Printx(x interface{}) {
	v := reflect.ValueOf(x)
	t := v.Type()
	fmt.Printf("type %s\n", t)
	// interface 有方法集, 才会有 NumMethod 方法 ???
	for i := 0; i < v.NumMethod(); i++ {
		methType := v.Method(i).Type()
		fmt.Printf("func (%s) %s%s\n", t, t.Method(i).Name,
			strings.TrimPrefix(methType.String(), "func"))
		// Trim 修剪

		/* TODO: 深入理解
		 *  reflect.Type 和 reflect.Value 都有方法 Method
		 *  reflect.Type 的 Method 方法会返回一个 reflect.Method 实例, 这个
		 *     结构类型描述了该方法的名称和类型
		 *
		 *  reflect.Value 的 Method 方法会返回一个 reflect.Value,
		 *	   代表一个方法值,  即一个已绑定接收者的方法
		 */
	}
}

func main() {
	Printx(time.Hour) // 一个类型作为参数???
}

/*
* type time.Duration
* func (time.Duration) Hours() float64
* func (time.Duration) Microseconds() int64
* func (time.Duration) Milliseconds() int64
* func (time.Duration) Minutes() float64
* func (time.Duration) Nanoseconds() int64
* func (time.Duration) Round(time.Duration) time.Duration
* func (time.Duration) Seconds() float64
* func (time.Duration) String() string
* func (time.Duration) Truncate(time.Duration) time.Duration
 */

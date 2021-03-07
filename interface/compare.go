package main

// 接口值的动态类型一致，但对应的动态值
// 如果是不可进行比较的, 进行接口值比较会崩溃(如 slice)
// 当把接口作为map的键或者switch语句的
// 操作数时，也存在崩溃的风险
// 仅在能确认接口值包含的动态值可
// 以比较的时候，才比较接口值

import "fmt"

func main() {
	var x interface{} = []int{1, 2, 3}
	fmt.Println(x == x)
	//  panic: 运行时错误吗，比较不能比较的类型 []int
}

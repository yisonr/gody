package main

/*
	defer 语句可以调试一个复杂的函数, 即在函数的"入口" 和 "出口" 处设置调试
	行为, 函数和参数表达式会在语句执行时求值, 无论是正常情况下执行return,
	或函数执行完毕, 还是不正常的情况下, 比如发生宕机, 实际的调用都推迟到
	包含defer语句的函数结束后才执行, defer 语句没有限制调用次数; 执行的时候
	以调用defer语句顺序的倒序进行(绑定到当前goroutine结构体, 入栈操作,
	TODO: 熟悉goroutine相关操作源码)
*/

import (
	"fmt"
	"log"
	"os"
	"time"
)

func bigSlowOperation() {
	defer trace("bigSlowOperation")() // 注意末尾括号的作用
	// ...
	time.Sleep(1 * time.Second) // 通过休眠仿真慢操作
}

func trace(msg string) func() {
	start := time.Now()
	return func() {
		log.Printf("exit %s (%s)", msg, time.Since(start))
	}
}

// func main() {
// 	bigSlowOperation()
// }

/*
	延迟执行的函数在 return 语句之后执行, 并且可以更新函数的结果变量,
 	因为匿名函数可以得到其外层函数作用域内的变量(包括命名的结果), 所以
 	延迟执行的匿名函数可以观察到函数的返回结果
*/

func double(x int) (result int) {
	defer func() { fmt.Printf("double(%d) = %d\n", x, result) }()
	return x + x
}

// func main() {
// 	_ = double(4) // "double(4) = 8"
// }

// 延迟执行的函数能改变外层函数返回给调用者的结果
func triple(x int) (result int) {
	defer func() { result += x }()
	return double(x)
}

// func main() {
// 	result := triple(4)
// 	fmt.Println(result)
// }

// 延迟函数不到外层函数的最后一刻不会执行, 所以应注意循环里 defer  语句的使用
func loopForDefer(filenames []string) error {
	for _, filename := range filenames {
		f, err := os.Open(filename)
		if err != nil {
			return err
		}

		defer f.Close() // 可能会用尽文件描述符, TODO: Linux文件描述符
	}
	// ...
	return nil
}

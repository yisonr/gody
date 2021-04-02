package main

import "fmt"

// version1
// 高度并行:由一些完全独立的子过程组成的过程,
// 因为子过程没有顺序上的相互关系，所以高度并行
// 是最容易实现的，有许多并行机制来实现线性扩展。
// 运行速度极快，甚至在文件名称slice中只有一个
// 元素的情况下，它也比不使用并发的情况快的多，
// 因为该程序在没有执行完任务之后就返回了。
// 主 goroutine 启动了所有 goroutine , 但是没有
// 等他们执行完毕。
// func main() {
// 	ch := make(chan struct{})
// 	filenames := []string{"cache.go", "test.py", "para.jpg"}
// 	for _, f := range filenames {
// 		go func() {
// 			fmt.Println(f)
// 			ch <- struct{}{}
// 		}()
// 	}
// 	// 没有输出，因为主 goroutine 直接暴力终结
// 	// 其他 goroutine, version2 中使用通道发送事件
// 	// 向外层 goroutine  报告它的完成。
// }

// version2
// 单变量f的值被所有匿名函数值共享并且被后续的
// 迭代所更新。新的 goroutine 执行字面量函数时，
// for 循环可能已经更新f,并且开始另一个迭代
// 或者迭代已经完全结束，所以当这些 goroutine
// 读取f的值时，取到的都是slice的最后一个元素。
// 换句话说，在for循环中使用并发(go关键字),
// 不会等待循环体的执行后才开始下一次循环而是会
// 直接开始循环;所以这里的输出都是一样的。
// func main() {
// 	ch := make(chan struct{})
// 	filenames := []string{"cache.go", "test.py", "para.jpg"}
// 	for _, f := range filenames {
// 		go func() {
// 			fmt.Println(f)
// 			ch <- struct{}{}
// 		}()
// 	}

// 	// 等待 goroutine 完成
// 	for range filenames {
// 		<-ch
// 	}
// }

// version3
// 字面量显式参数传递f,而不是在 for 循环中声明，
// 单变量f的值被所有匿名函数值共享，那么非匿名函数???
func main() {
	ch := make(chan struct{})
	filenames := []string{"cache.go", "test.py", "para.jpg"}
	for _, f := range filenames {
		go func(x string) {
			fmt.Println(x)
			ch <- struct{}{}
		}(f)
	}

	// 等待 goroutine 完成
	for range filenames {
		<-ch
	}
}

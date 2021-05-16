package main

import (
	"fmt"
	"sync"
	"time"
)

// 直觉来看，可能觉得 loadIcons() 会执行多次,
// 但所有关于并发的直觉，都是不可靠的
//
// 在缺乏显式同步的情况下，编译器和cpu在能保证每个
// goroutine 都满足串行一致性的基础上可以自由重排
// 访问内存的顺序.
//
// 注意一个 goroutine 发现 icons 不是 nil 并不意味着
// 变量的初始化肯定已经完成

// 在实际的 goroutine 调度中，B 先判断出 icons 为 nil,
// 然后便开始调用 loadIcons(), 在调用到 32 行后，调度器
// 开始调度 A goroutine, 此时icons不为nil, 且长度为 1,
// 此时 A 取 incos[clubs.png] 时键值不存在，得到字符串零值
// A 调度完后，B 继续从 33 行执行，直至 goroutine 串行完成。

var icons map[string]string

func loadIcon(name string) string {
	return name
}

func loadIcons() {
	icons = map[string]string{
		"spades.png": loadIcon("spades.png"),
	}
	// B 退出调度，A 进入调度，此时 icons 的长度为1
	// A 中访问的 clubs.png 键值不存在
	time.Sleep(100 * time.Microsecond)
	// 延迟，防止执行过快, 未用完调度时间片
	// 不能观察调度切换现象
	icons = map[string]string{
		"hearts.png":   loadIcon("hearts.png"),
		"diamonds.png": loadIcon("diamonds.png"),
		"clubs.png":    loadIcon("clubs.png"),
	}
}

// 并发不安全
func Icon(P, name string) string {
	if icons == nil {
		fmt.Println(P + " is nil")
		loadIcons() // 一次性的初始化
	}
	return icons[name]
}

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() { // A goroutine
		defer wg.Done()
		fmt.Println(Icon("A", "clubs.png"))
		fmt.Println(len(icons))
		fmt.Println("---------")
	}()
	go func() { // B goroutine
		defer wg.Done()
		fmt.Println(Icon("B", "clubs.png"))
		fmt.Println(len(icons))
		fmt.Println("/////////")
	}()
	wg.Wait()
}

// 预测结果:
//B is nil
//
//1
//--------
//clubs.png
//3
//////////

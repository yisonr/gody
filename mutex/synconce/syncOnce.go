package main

// ./iconMutex.go 的改进版

import (
	"fmt"
	"sync"
	"time"
)

// sync 包提供的一次性初始化问题的特化解决方案:
// sync.Once
// Once 包含一个布尔变量和一个互斥量，布尔变量记录
// 初始化是否已经完成，互斥量则负责保护这个布尔变量和
// 客户端的数据结构，Once  的唯一方法Do以初始化函数
// 作为它的参数.

// 每次调用 Do(loadIcons) 时会先锁定互斥量并检查里边的
// 布尔变量，若这个布尔变量为假，Do 会调用 loadIcons
// 然后把变量设置为真，后续的调用相当于空操作，只是通过
// 互斥量的同步来保证 loadIcons 对内存产生的效果，即 icons
// 变量对所有的 goroutine 可见，以这种方式来使用 sync.Once
// 可以避免变量在正确构造之前就被其他的 goroutine 共享.

// todo: 查看 sync.Once 源码，研究其实现原理

var icons map[string]string

func loadIcon(name string) string {
	return name
}

func loadIcons() {
	icons = map[string]string{
		"spades.png": loadIcon("spades.png"),
	}
	time.Sleep(100 * time.Microsecond)
	icons = map[string]string{
		"hearts.png":   loadIcon("hearts.png"),
		"diamonds.png": loadIcon("diamonds.png"),
		"clubs.png":    loadIcon("clubs.png"),
	}
}

var loadIconOnce sync.Once

// 并发安全
func Icon(P, name string) string {
	loadIconOnce.Do(loadIcons)
	return icons[name]
	// if icons == nil {
	// 	fmt.Println(P + " is nil")
	// 	loadIcons() // 一次性的初始化
	// }
	// return icons[name]
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

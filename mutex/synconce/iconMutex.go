package main

// ./icon.go 的并发安全版

import (
	"fmt"
	"sync"
	"time"
)

var icons map[string]string

func loadIcon(name string) string {
	return name
}

func loadIcons() {
	icons = map[string]string{
		"diamonds.png": loadIcon("diamonds.png"),
	}
	time.Sleep(100 * time.Microsecond)
	icons = map[string]string{
		"hearts.png":   loadIcon("hearts.png"),
		"diamonds.png": loadIcon("diamonds.png"),
		"clubs.png":    loadIcon("clubs.png"),
	}
}

var mu sync.RWMutex // 保护 icons

// 并发安全
func Icon(P, name string) string {
	// -------------------------------- 1 号临界区(可以多个goroutine读)
	mu.RLock() // 获取共享锁(读锁)
	if icons != nil {
		icon := icons[name]
		mu.RUnlock()
		return icon
	}
	mu.RUnlock()
	// -------------------------------- 1 号临界区
	// 由于不先释放一个共享锁就无法直接把它升级为
	// 互斥锁，为了避免在过渡期其他 goroutine  已经
	// 初始化了 icons, 所以必须重新检查 nil 值
	// ******************************** 2 号临界区(未找到条目后开始互斥初始化)
	mu.Lock() // 获取互斥锁(写锁)

	if icons == nil { // 必须重新检查nil的值
		fmt.Println(P + " is nil")
		loadIcons() // 一次性的初始化
	}
	icon := icons[name]
	mu.Unlock()
	// ******************************** 2 号临界区
	return icon
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

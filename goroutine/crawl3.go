package main

// TODO: 并发套用的作用，31行不并发为什么会出现死锁
// 练习8.6
// 练习8.7

import (
	"fmt"
	"log"
	"os"
	"time"

	"gopl.io/ch5/links"
)

func crawl(url string) []string {
	fmt.Println(url)
	list, err := links.Extract(url)
	if err != nil {
		log.Print(err)
	}
	return list
}

func main() {
	worklist := make(chan []string)  // 可能有重复的url列表
	unseenLinks := make(chan string) // 去重后的url列表

	// 从命令行参数开始, 发送给任务列表的命令行参数
	// 必须在它自己的goroutine中允许来避免死锁???
	go func() { worklist <- os.Args[1:] }()

	// 创建20个爬虫goroutine来获取每个可见链接
	for i := 0; i < 20; i++ {
		go func() {
			for link := range unseenLinks {
				foundLinks := crawl(link)
				go func() { worklist <- foundLinks }() // 为什么要并发???
			}
		}()
	}

	// 主 goroutine 对url列表进行去重
	// 并把没有爬取过的条目发送给爬虫程序
	seen := make(map[string]bool)
	for list := range worklist {
		for _, link := range list {
			if !seen[link] {
				seen[link] = true
				unseenLinks <- link
			}
		}
	}
}

package main

import (
	"fmt"
	"log"
	"os"

	"gopl.io/ch5/links"
)

// 令牌是一个计数信号量
// 确保并发请求限制在20个以内
var tokens = make(chan struct{}, 20)

func crawl(url string) []string {
	fmt.Println(url)
	tokens <- struct{}{}            // 获取令牌
	list, err := links.Extract(url) // 临界区???
	<-tokens                        // 释放令牌
	if err != nil {
		log.Print(err)
	}
	return list
}

func main() {
	worklist := make(chan []string)
	var n int // 等待发送到任务列表的数量
	// 计数器n跟踪发送到任务列表中的任务个数,每次知道一个条目被发送
	// 到任务列表时，就递增变量n，第一次递增是在发送初始化命令行参数之前
	// 第二次递增是在每次启动一个新的爬取goroutine的时候。主循环从n减到0,
	// 此时再没有任务需要完成。

	// 从命令行参数开始
	n++ // ???
	go func() { worklist <- os.Args[1:] }()

	// 并发爬取 web
	seen := make(map[string]bool)
	for ; n > 0; n-- {
		list := <-worklist
		for _, link := range list {
			if !seen[link] {
				seen[link] = true
				n++
				go func(link string) {
					worklist <- crawl(link)
				}(link)
				// 字面量函数的显式参数传递, 避免循环变量捕获问题
			}
		}
	}
}

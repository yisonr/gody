package main

import (
	"fmt"
	"log"
	"os"

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
	worklist := make(chan []string)

	// 从命令行参数开始
	go func() { worklist <- os.Args[1:] }()

	// 并发爬取 web
	seen := make(map[string]bool)
	for list := range worklist {
		for _, link := range list {
			if !seen[link] {
				seen[link] = true
				go func(link string) {
					worklist <- crawl(link)
				}(link)
				// 字面量函数的显式参数传递, 避免循环变量捕获问题
			}
		}
	}
}

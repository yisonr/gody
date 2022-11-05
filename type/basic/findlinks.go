package main

/* 匿名函数的使用 */

import (
	"fmt"
	"log"
	"os"

	"gopl.io/ch5/links"
)

func main() {
	// 开始广度遍历
	breadthFirst(crawl, os.Args[1:])
}

// 广度优先遍历
// 将函数的行为当做参数传递
func breadthFirst(f func(item string) []string, worklist []string) {
	seen := make(map[string]bool)
	for len(worklist) > 0 {
		items := worklist // 维护一个字符串集合, 保证每个节点只访问一次
		worklist = nil
		for _, item := range items {
			if !seen[item] {
				seen[item] = true
				worklist = append(worklist, f(item)...)
			}
		}
	}
}

func crawl(url string) []string {
	fmt.Println(url)
	list, err := links.Extract(url)
	if err != nil {
		log.Println(err)
	}
	return list
}

// func Extract(url string) ([]string, error) {
// 	links.Extract(url string)
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if resp.StatusCode != http.StatusOK {
// 		resp.Body.Close()
// 		return nil, fmt.Errorf("getting %s: %s", url, resp.Status)
// 	}

// 	doc, err := html.Parse(resp.Body)
// 	resp.Body.Close()
// 	if err != nil {
// 		return nil, fmt.Errorf("parsing %s as HTML: %v", url, err)
// 	}

// 	var links []string
// 	visitNode := func(n *html.Node) {
// 		if n.Type == html.ElementNode && n.Data == "a" {
// 			for _, a := range n.Attr {
// 				if a.Key != "href" {
// 					continue
// 				}
// 				link, err := resp.Request.URL.Parse(a.Val)
// 				if err != nil {
// 					continue // 忽略不合法的URL
// 				}
// 				links = append(links, link.String())
// 			}
// 		}
// 	}
// 	forEachNode(doc, visitNode, nil)
// 	return links, nil
// }

package main

/*
	递归
	"golang.org/x/net/html" 提供了解析html的功能
	许多编程语言使用固定长度的函数调用栈: 大小在64KB到2MB之间; 递归的深度会
	受限于固定长度的栈大小, 所以当进行深度递归使用时必须谨防栈溢出; 固定长度
	的栈甚至会造成一定的安全隐患; 相比固定长度的栈, go的实现使用了可变长度的
	栈, 栈的大小会随着使用而增长, 可达到1GB左右的上限, 所以可以安全的使用递归
	而不用担心溢出问题;
	wget https://golang.org
	cat index.html | go run crawl.go

*/

import (
	"fmt"
	"os"

	"golang.org/x/net/html"
)

func main() {
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "findlinks1: %v\n", err)
		os.Exit(1)
	}
	// for _, link := range visit(nil, doc) {
	// 	fmt.Println(link)
	// }

	outline(nil, doc)
}

/*
	visit 递归的遍历 HTML 树上的所有节点, 从 HTML 锚元素 <a href='...'> 中得到
	href 属性的内容, 将获取到的链接内容添加到字符串 slice 并返回;

	要对树中的任意节点 n 进行递归, visit 递归的调用自己去访问节点n的所有子节点,
	并且将访问过的节点保存在 FirstChild 链表中;
*/
func visit(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				links = append(links, a.Val)
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = visit(links, c)
	}
	return links
}

// 使用递归遍历所有html文本中的节点树, 并输出树的结构, 当递归遇到每个元素时,
// 都会将元素标签压入栈, 然后输出栈
func outline(stack []string, n *html.Node) {
	if n.Type == html.ElementNode {
		// stack := append(stack, n.Data) // 注意这种写法会造成的错误
		stack = append(stack, n.Data) // 把标签压入栈
		fmt.Println(stack)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		outline(stack, c)
	}
}

package main

/* 匿名函数的使用 */

import (
	"fmt"
	"sort"
)

// 反映了所有课程和先决课程的关系
// 考虑学习计算机课程的顺序, 需要计算出学习每一门课程的先决条件
// 先决条件的内容构成了一张有向图, 每个节点代表每一门课程, 每一条边代表一门
// 课程所依赖另一门课程的关系, 图是无环的; 没有节点可以通过图上的路径回到它自己
var prereqs = map[string][]string{
	"algorithma": {"data structures"},
	"calculus":   {"linear algebara"},
	"compiles": {
		"data structures",
		"formal languages",
		"computer organization",
	},
	"data structures":       {"discrete math"},
	"databases":             {"data structures"},
	"discrete math":         {"intro to programming"},
	"formal languages":      {"discrete math"},
	"networks":              {"operating systems"},
	"operating systems":     {"data structures", "computer organziation"},
	"programming languages": {"data structures", "computer organziation"},
}

func main() {
	for i, course := range topoSort(prereqs) {
		fmt.Printf("%d:\t%s\n", i+1, course)
	}
}

// 拓扑排序, TODO
// 使用深度优先的搜索算法计算得到合法的学习路径
func topoSort(m map[string][]string) []string {
	var order []string
	seen := make(map[string]bool)

	// 当一个匿名函数需要进行递归, 必须先声明一个变量然后将匿名函数赋给这个
	// 变量; 如果将两个步骤合并成一个声明, 函数字面量将不能存在于 visitAll
	// 变量的作用域中, 这样就不能递归的调用自己了
	var visitAll func(items []string)
	visitAll = func(items []string) {
		for _, item := range items {
			if !seen[item] {
				seen[item] = true
				visitAll(m[item])
				order = append(order, item)
			}
		}
	}

	var keys []string
	for key := range m {
		keys = append(keys, key)
	}

	sort.Strings(keys)
	visitAll(keys)
	return order
}

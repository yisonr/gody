package main

/* map 是一个拥有键值对元素的无序集合, 其键必须是可以比较的数据类型 */

import (
	"bytes"
	"fmt"
	"sort"
)

//
// 映射: map
// 可以使用类似处理数组和切片的方式迭代映射中的元素, 但映射是无序的， 即
// 每次迭代映射的时候顺序都是不一致的
// 无序的原因是映射的实现使用了散列表
// 映射的散列表包含一组桶，在存储、删除或查找键值对的时候，所有的操作都要
// 先选择一个桶, 把操作映射时指定的键传给映射的散列函数，就能选中对应的桶
// 这个散列函数的目的是生成一个索引，最终将键值对分布到所有可用的桶里
//
//
// todo:
// 随着映射存储的增加，索引分布越均匀，访问键值对的速度就越快???
// 映射通过合理数量的桶来平衡键值对的分布
// 对 go 语言的映射来说，生成的散列键的一部分，具体来说是低位(LOB),
// 被用来选择桶

func main1() {
	// map 是引用类型, 当传递map给一个函数, 并对这个映射做了修改时, 所有对
	// 这个映射的引用都会感知到这个修改

	colors := map[int]string{
		1: "测试1",
		2: "测试2",
		3: "测试3",
		4: "测试4",
	}

	for key, value := range colors {
		fmt.Printf("Key: %v value: %s\n", key, value)
	}

	fmt.Println("----------------------------")
	removeColor(colors, 1)
	for key, value := range colors {
		fmt.Printf("Key: %v value: %s\n", key, value)
	}
}

func removeColor(colors map[int]string, key int) {
	delete(colors, key)
}

//
// Key: 2 value: 测试2
// Key: 3 value: 测试3
// Key: 4 value: 测试4
// Key: 1 value: 测试1
// ----------------------------
// Key: 3 value: 测试3
// Key: 4 value: 测试4
// Key: 2 value: 测试2
//

/*
   map 的元素不是一个变量, 不能获取其地址
   _ = &ages["bob"]  // 编译错误
   不能获取 map 元素地址的原因是 map 的增长可能会导致已有元素被重新散列到新的
   存储位置, 这样就可能会使获取的地址无效;(TODO: map的增长原理)

   这就导致map的value为结构体时更新value时不能用 map[key].Field 的方式
   https://segmentfault.com/q/1010000041583173
   https://go.dev/ref/spec#Assignments
   The operand must be addressable, that is, either a variable,
   pointer indirection, or slice indexing operation; or a field selector
   of an addressable struct operand; or an array indexing operation
   of an addressable array.
*/

// 如果要给map排序, 必须显式的给键排序; 如果键是字符串类型, 可以使用 sort 包
// 中的 Strings 函数给键排序
func sortMapKey(ages map[string]int) {
	// var names []string
	names := make([]string, 0, len(ages))
	// 创建一个初始元素为空, 但容量足够容纳 ages map 中所有键的 slice, 相比
	// []string 的方式更加高效(需要扩容)
	for name := range ages {
		names = append(names, name)
	}

	sort.Strings(names)
	for _, name := range names {
		fmt.Printf("%s\t%d\n", name, ages[name])
	}
}

/*
	查找元素, 删除元素, 获取map元素个数, 执行 range 循环都可以在 map 的零值 nil
	上执行, 但是向零值 map 中设置元素会 panic, 设置元素之前, 必须初始化 map

	和 slice 一样, map 不可比较, 唯一合法比较就是和 nil 比较


	有时候需要一个 map 并且它的键是 slice, 但是 map 的键必须是可比较的, 所以需
	要使用下列思路实现:
	- 定义一个辅助函数k将每个键都映射到字符串, 仅当x和y相等的时候, 才认为 k(x) == k(y)
	- 创建一个map, map的键是字符串类型, 在每个键元素被访问的时候, 调用这个帮助函数
*/
// 通过一个字符串列表使用一个map来记录Add函数被调用的次数
var m = make(map[string]int)

// 使用 %q 格式化 slice 并记录每个字符串的边界
// %q: 双引号围绕的字符串, 由Go语法安全地转义
func k(list []string) string  { return fmt.Sprintf("%q", list) }
func Add(list []string)       { m[k(list)]++ }
func Count(list []string) int { return m[k(list)] }

func main2() {
	var data1 = []string{"sd", "bo", "oisud"}
	Add(data1)
	fmt.Println(k(data1))
	fmt.Println(Count(data1))
}

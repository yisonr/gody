package main

/* slice 就地修改 */

// noempty 从给定的字符串列表中去除空字符串并返回新的slice
func noempty(strings []string) []string {
	// 输入的 slice 和输出的 slice 拥有相同的底层数组, 避免在函数内重新分配数组
	i := 0
	for _, s := range strings {
		if s != "" {
			strings[i] = s
			i++
		}
	}
	return strings[:i]
}

/*
	无论使用哪种方式, 重用底层数组的结果是每一个输入值的 slice 最多只有一个
	输出的结果 slice, 很多从序列中过滤元素再组合结果的算法都是如此

	// slice 实现栈
	stack = append(stack, v) // push v
	top := stack[len(stack)-1] // 栈顶
	stack = stack[:len(stack)-1] // pop
*/

//@4 移除slice元素的高级实现
func remove(s []int, i int) []int {
	// 使用 copy 将高位索引的元素向前移动来覆盖被移除元素所在位置
	copy(s[i:], s[i+1:])
	return s[:len(s)-1]
}

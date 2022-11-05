package main

import "os"

/*
	警告: 捕获迭代变量
	TODO: 词法作用域
*/

func mkdirAndRm(tempDirs func() []string) {
	var rmdirs []func()
	for _, d := range tempDirs() {
		dir := d // 此行是必须的, 在循环里创建的所有函数变量共享
		// 相同的变量--一个可访问的存储位置, 而不是固定的值, 当调用
		// 清理函数时, d 变量已经被每一次的for循环更新多次, 因此引入
		// 内部变量 dir 来解决此问题
		// dir := dir // 声明内部的dir, 并以外部的dir初始化
		os.MkdirAll(dir, 0755) // 同时创建父目录
		rmdirs = append(rmdirs, func() {
			os.RemoveAll(dir)
		})
	}

	// ...

	for _, rmdir := range rmdirs {
		rmdir()
	}
}

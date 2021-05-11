package main

// go 程序语言设计
// Page 209
// 同步不仅涉及多个 goroutine 的执行顺序问题，
// 还会影响到内存
// todo: 实验
// ./test.go 出现奇怪的结果

import (
	"fmt"
	"sync"
)

func Memory() string {
	var wg sync.WaitGroup
	exam := ""
	wg.Add(2)
	var x, y int
	go func() {
		defer wg.Done()
		x = 1
		exam = exam + fmt.Sprintf("y: %v ", y)
	}()
	go func() {
		defer wg.Done()
		y = 1
		exam = exam + fmt.Sprintf("x: %v ", x)
	}()
	wg.Wait()
	return exam
}

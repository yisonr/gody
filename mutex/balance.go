package main

var balance int

func Deposit(amount int) {
	balance = balance + amount
}

func Balance() int {
	return balance
}

/*
对上面两个函数进行并发调用, 会出现数据竞态
数据竞态发生于两个 goroutine 并发读写同一个变量
并且至少其中一个是写入时

*/

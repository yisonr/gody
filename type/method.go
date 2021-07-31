package main

import "fmt"

/*
 * 方法能为用户定义的类型添加新的行为
 * 方法实际上也是函数，只是这个函数与类型绑定，只有方法的接收者才可调用方法
 * 关键字 func 和函数名之间的参数被称为接收者, 将函数和接收者的类型绑定到
 * 一起，有接收者的函数称为方法
 *
 * go 语言里有两种类型的接收者: 值接收者和指针接收者
 *
 */

type user struct {
	name  string
	email string
}

// notify 使用值接收者实现了一个方法
func (u user) notify() {
	fmt.Printf("Sending User Email To  %s<%s>\n", u.name, u.email)
}

// 如果使用值接收者声明方法， 调用时会使用这个值的一个副本来执行.

// 也可以使用指针调用使用值接收者声明的方法
// lisa := &user{"Datou", "datou@email.com"}
// lisa.notify()
// 为了支持这种方法调用, go 语言调整了指针的值，来符合方法接收者的定义
// (*lisa).notify(), 指针被解引用为值，符合了值接收者的要求，注意: notify
// 操作的是一个副本，只不过这次操作的是从 lisa 指针指向的值的副本.

// changeEmail 使用指针接收者实现了一个方法
// changeEmail 的调用对值做的修改会反映在 lisa 指针所指向的值上,
func (u *user) changeEmail(email string) {
	u.email = email
}

// 值接收者使用值的副本来调用方法，而指针接收者使用实际值来调用方法,
// 也可以使用一个值来调用使用指针接收者声明的方法
// bill := user{"xiaohong", "xiaohong@email.com"}
// bill.changeEmail("xiaohong@qq.com")
// go 语言对值做了调整，使之符合函数的接收者并进行调用
// (&bill).changeEmail("xiaohong@qq.com")

// go 语言既允许使用值，也允许使用指针来调用方法，不必严格符合接收者的类型

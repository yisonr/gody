package main

// 此节配合 ./interface.go 理解

// 多态: polymorphism
/*
 * 多态:
 * - 指为不同数据类型的实体提供统一的接口，多态类型可以将自身
 *  所支持的操作套用到其他类型的值上
 * - 指代码可以根据类型的具体实现采取不同行为的能力
 *
 *
 * 接口是用来定义行为的类型，这些被定义的行为不由接口直接实现，而是通过方法由
 * 用户定义的类型实现。如果用户定义的类型实现了某个接口类型声明的一组方法，
 * 那么这个用户定义的类型的值就可以赋给这个接口类型的值, 这个赋值会把用户定义
 * 的类型的值存入接口类型的值
 *
 *
 * 对接口值方法的调用会执行接口值里存储的用户定义的类型的值的对应的方法，
 * 任何用户定义的类型都可以实现任何接口，所以对接口值方法的调用自然就是
 * 一种多态, 在这个关系里，用户定义的类型通常叫做实体类型，原因是如果离开
 * 内部存储的用户定义的类型的值的实现，接口值并没有具体的行为
 *
 */

// 如下例子:

import "fmt"

// 定义了通知类行为的接口
type notifier interface {
	notify()
}

// 定义的用户类型
type user struct {
	name  string
	emial string
}

func (u *user) notify() {
	fmt.Printf("Sending user email to %s<%s>\n", u.name, u.emial)
}

// 定义了管理员
type admin struct {
	name  string
	email string
}

func (a *admin) notify() {
	fmt.Printf("Sending admin email to %s<%s>\n", a.name, a.email)
}

func main() {
	bill := user{"Bill", "bill@email.com"}
	sendNotification(&bill)
	// *user 接口值的方法集为: (u user) 和 (u *user)

	lisa := admin{"Lisa", "lisa@email.com"}
	sendNotification(&lisa)
}

// SendNotification 接受一个实现了 notifier 接口的值, 并发送通知
func sendNotification(n notifier) {
	n.notify()
}

// 任意一个实体类型都能实现 notifier 接口，那么函数 sendNotification 可以
// 针对任意实体类型的值来执行 notifier 方法，因此，函数 sendNotification
// 就能提供多态的行为

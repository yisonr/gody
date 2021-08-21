package main

// 外部类型重写内部类型的实现

import "fmt"

type notifier interface {
	notify()
}

type user struct {
	name  string
	email string
}

// 通过 user 类型值的指针调用的方法
func (u *user) notify() {
	fmt.Printf("Sending user email to %s<%s>\n",
		u.name,
		u.email)
}

// 一个拥有权限的管理员用户
type admin struct { // user 是外部类型 admin 的内部类型
	user  // 嵌入类型
	level string
}

// 通过 admin 类型值的指针调用的方法
func (a *admin) notify() {
	fmt.Printf("Sending admin email to %s<%s>\n",
		a.name,
		a.email)
}

func main() {
	// 创建一个 admin 用户
	ad := admin{
		user: user{
			name:  "loin",
			email: "loin@163.com",
		},
		level: "super",
	}

	// 给admin用户发送一个通知
	// 接口的嵌入的内部类型实现并没有提升到外部类型
	sendNotification(&ad)

	// 访问内部类型的方法
	ad.user.notify()

	// 访问外部类型的方法(内部类型的方法没有被提升,或理解为被覆盖)
	ad.notify()

	/*
			 * 若外部类型实现了 notify 方法，内部类型的实现就不会被提升，
		     * 内部类型的值一直存在，因此还可以通过直接访问内部类型的值来
			 * 调用没有被提升的内部类型的实现的方法
	*/
}

func sendNotification(n notifier) {
	n.notify()
}

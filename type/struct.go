package main

// 嵌入类型
// go 语言允许用户扩展或修改已有类型的行为, 嵌入类型是将已有的
// 类型直接声明在新的结构类型里, 被嵌入的类型称为外部类型的
// 内部类型

// 通过嵌入类型，与内部类型相关的标识符会提升到外部类型上，被提升
// 的标识符与直接声明在外部类型里的标识符一样，也是外部类型的一部分,
// 这样外部类型就组合了内部类型包含的所有属性，并且可以添加新的字段和
// 方法，外部类型也可以通过声明与内部类型标识符同名的标识符来覆盖内部
// 标示符的字段或方法，这即扩展或修改已有类型的方法.

import "fmt"

type user struct {
	name  string
	email string
}

type notifier interface {
	notify()
}

func (u *user) notify() {
	fmt.Printf("Sending user email to %s<%s>\n",
		u.name,
		u.email)
}

type admin struct { // user 是外部类型 admin 的内部类型
	user  // 嵌入类型
	level string
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

	// 直接访问内部类型的方法
	// 对外部类型来说,内部类型总是存在的, 即虽然没有指定内部类型对应
	// 的字段名, 还是可以使用内部类型的类型名来访问内部类型的值?
	// ad.user.notify()

	// 内部类型的方法被提升到外部类型
	// ad.notify()
	// 借助内部类型提升，notify 方法可以直接通过ad变量来访问

	// 给admin用户发送一个通知
	// 用于实现接口的内部类型的方法，被提升到外部类型
	// 注意方法集的理解
	sendNotification(&ad)
	// 将外部类型的地址传给 sendNotification 函数, 编译器认为
	// 这个指针实现了 notifier 接口, 并接受了这个值的传递, 但
	// admin 类型并没有实现这个接口。
	// 由于内部类型的提升，内部类型实现的接口会自动提升到外部类型,
	// 即由于内部类型的实现, 外部类型也同样实现了这个接口

}

func sendNotification(n notifier) {
	n.notify()
}

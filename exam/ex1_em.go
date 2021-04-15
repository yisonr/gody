package main

//  todo:
// 内存的分配
// for range 循环的内层逻辑

import "fmt"

type student struct {
	Name string
	Age  int
}

func main() {
	m := make(map[string]*student)
	stus := []student{
		{
			Name: "zhou",
			Age:  24,
		}, {
			Name: "li",
			Age:  12,
		},
	}
	for _, stu := range stus {
		fmt.Println(&stu)
		// fmt.Println(&stu)
		// m[stu.Name] = &student{Name: stu.Name, Age: stu.Age}
		m[stu.Name] = &stu // 这里取不到地址, 只能取到一个临时地址?
	}

	for _, stu := range stus {
		stu.Age = stu.Age + 10
	}

	for _, value := range m {
		fmt.Println(value.Name)
		fmt.Println(value.Age)
	}

}

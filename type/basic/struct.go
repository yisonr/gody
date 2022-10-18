package main

/*
	结构体的每个成员都可以通过点号访问, 通过点号访问的成员都是变量(即可以地址),
	因此可以给结构体成员赋值;

	TODO: 空结构体(https://zhuanlan.zhihu.com/p/351176221)
*/

type Employee struct {
	ID      int
	Name    string
	Address string
	Salary  int
}

func EmployeeByID(id int) *Employee {
	return &Employee{}
}

// func EmployeeByID(id int) Employee {
// 	return Employee{}
// }

func SetEmployee() {
	// EmployeeByID 函数的返回值是 Employee 类型而非 *Employee 类型时, 代码
	// 将无法编译, 因为赋值表达式的左侧无法识别出一个变量(不能取地址)

	// 函数此时的返回值 Employee 类型变量的作用域在 EmployeeByID 在函数
	//  内部, 函数的调用结果没有在调用它的函数内部进行显式传递, 又因为返回
	// 值类型不是指针类型, 所以没有发生内存逃逸确定变量存放在堆或栈上,
	// 即调用该函数的函数不能识别出一个变量
	// TODO: 内存逃逸
	EmployeeByID(1).Salary = 0
}

// 结构体中的成员变量通常一行写一个, 相同类型的连续成员变量可以写在一行上,
// 成员变量的顺序对于结构体同一性很重要(内存对齐), 如果成员顺序不一样则是
// 不同的结构体

/*
	go 最主要的访问控制机制: 如果结构体的成员变量名称是首字母大写的, 则此
	变量是可导出的;
	一个结构体可以同时包含可导出和不可导出的成员变量;


	命名结构体类型S不可以定义一个拥有相同结构体类型S的成员变量, 也就是一个
	聚合类型不可以包含它自己(同样的限制对数组也适用), 但是S中可以定义一个S的
	指针类型*S, 这样就可以创建一些递归数据结构, 比如链表和树;
*/

// 利用二叉树实现插入排序的例子
type tree struct {
	value       int
	left, right *tree
}

func Sort(values []int) {
	var root *tree
	for _, v := range values {
		root = add(root, v)
	}
	appendValues(values[:0], root)
}

func appendValues(values []int, t *tree) []int {
	if t != nil {
		values = appendValues(values, t.left)
		values = append(values, t.value)
		values = appendValues(values, t.right)
	}
	return values
}

func add(t *tree, value int) *tree {
	if t == nil {
		//  等价于返回 &tree{value: value}
		t = new(tree)
		t.value = value
		return t
	}

	if value < t.value {
		t.left = add(t.left, value)
	} else {
		t.right = add(t.right, value)
	}
	return t
}

/*
	结构体零值有结构体成员的零值组成, 通常, 希望零值是一个默认自然的, 合理的
	值; 如 bytes.Buffer 中, 结构体的初始值就是一个可以直接使用的空缓存;
	sync.Mutex 也是一个可以直接使用并且未锁定状态的互斥锁; 有时候, 这种合理的
	初始值实现简单, 但是有时候也需要类型的设计者花费时间进行设计;

	没有任何成员变量的结构体称为空结构体(struct{}), 没有长度也不携带任何信息,
	但是有时候会很有用, 可以用来替代被当做集合使用的map中的布尔值来强调只有
	键是有用的;

	如果结构体的所有变量都可以比较, 那么这个结构体就是可比较的; 两个结构体
	的比较可以使用 == 或者 !=, 其中 == 操作符按照顺序比较两个结构体变量的成
	员变量; 和其他可比较类型一样, 可比较的结构体类型都可以作为map的键类型;
*/

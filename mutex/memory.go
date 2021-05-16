package main

// go 程序语言设计
// Page 209
// 同步不仅涉及多个 goroutine 的执行顺序问题，
// 还会影响到内存
// 现代的计算机一般会有多个处理器，每个处理器都有内存的
// 本地缓存。为了提高效率，对内存的写入是缓存在每个处理器中
// 的，只有在必要时才刷回内存。甚至刷回内存的顺序都可能与
// goroutine 的写入顺序不一致，像通道通信或者互斥锁操作这样
// 的同步原语都会导致处理器把累积的写操作刷回内存并提交，
// 所以这个时刻之前 goroutine 的执行结果就保证了对运行在
// 其他处理器的  goroutine 可见。
// todo: 实验
// ./memoryTest.go 出现奇怪的结果

// 在单个 goroutine 内, 每个语句的效果保证按照执行的顺序发生，
// 即 goroutine 是串行一致的。但是在缺乏使用通道或者互斥量来显
// 式同步的情况下， 并不能保证所有的 goroutine 看到的事件顺序都
// 是一致的，尽管 goroutine A 肯定能在读取y之前观察到 x=1 的效果，
// 但它并不一定能观察到 goroutine B 对y写入的效果，所以A可能会
// 输出y的一个过期值。
// 尽管很容易把并发简单理解为多个 goroutine 中语句的某种交错执行
// 方式，但这并不是一个现代编译器和cpu的工作方式，因为赋值和 print
// 对应对应不同的变量，所有编译器就会认为两个语句的执行顺序不会影响
// 结果，然后就交换了这个语句的执行顺序，CPU也有类似的问题，如果两个
// goroutine 在不同的cpu上执行，每个cpu都有自己的缓存，那么一个goroutine
// 的写入操作在同步到内存之前对另外一个goroutine的print语句是不可见的。
//
// 这些并发问题都可以通过简单的、成熟的模式来避免，即在可能的情况下，
// 把变量限制到单个 goroutine 中，对于其他变量，使用互斥锁。
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
		x = 1                                  // A1
		exam = exam + fmt.Sprintf("y: %v ", y) // A2
	}()
	go func() {
		defer wg.Done()
		y = 1                                  // B1
		exam = exam + fmt.Sprintf("x: %v ", x) // B2
	}()
	wg.Wait()
	return exam
}

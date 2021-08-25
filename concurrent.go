package main

// go 语言的并发指的是能让某个函数独立于其他函数运行的能力,
// 当一个函数创建为 goroutine 时, go 会将其作为一个独立的工作单元,
// 这个单元会被调度到逻辑处理器上执行。 go 语言的调度器运行在
// 操作系统之上, 将操作系统的线程与语言运行时的逻辑处理器绑定，并
// 在逻辑处理器上运行 goroutine, 调度器在任何给定的时间，都会全面
// 控制哪个 goroutine 在哪个逻辑处理器上运行。

// TODO: 深入了解
// 泛型CSP(Communicating Sequential Process), 通信顺序进程
// csp 是一种消息传递模型，通过在 goroutine 之间传递数据来传递
// 消息，而不是对数据进行加锁来实现同步访问.
// 用于在 goroutine 之间同步或传递数据的关键数据模型叫做通道(channel)

/*
* TODO: 百度所言，持怀疑态度
* CSP 作为理论支持, go 实现对并发的原生支持，从实际出发，go并未完全
* 实现 csp 模型的所有理论，仅是借用 process 和 channel 这两个概念，
* process 在 go 语言上的实现就是 goroutine, 是实际并发执行的实体，
* 每个实体间通过 channel 通讯来实现数据共享。
*
 */

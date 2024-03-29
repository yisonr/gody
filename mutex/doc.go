package main

/*
并发安全
如果一个函数在并发调用时仍然能正确工作，那么此函数是并发安全的
并发调用是指: 在没有额外同步机制的情况下，从两个或多个 goroutine
同时调用这个函数。可以推广到方法或者作用于特定类型的一些操作.
如果一个类型的所有可访问方法和操作都是并发安全的，则称为并发安全的类型。

仅在文档指出类型是安全的情况下，才可以并发的访问一个变量，对于绝大部分变量，
如要回避并发访问，要么限制变量只存在于一个 goroutine 内，要么维护一个更高层的
互斥不变量。


因为包级别的变量无法限制在一个 goroutine  内，修改这些变量的函数就必须
采用互斥机制，所以导出的包级别的函数通常可以认为是并发安全的。


函数并发调用时出现问题: 死锁(deadlock)，活锁(livelock)，资源耗尽
活锁: 多个线程在尝试绕开死锁，却由于过分同步导致反复冲突

竞态:
指多个 goroutine 按某些交错顺序执行时程序无法给出正确的结果.
竞态潜伏在程序中，出现频率低，有可能仅在高负载环境或者特定的编译
器、平台和架构时才出现。使得竞态难以再现和分析

数据竞态是竞态中的一种，称为数据竞态，数据竞态发生于两个 goroutine
并发读写同一个变量并且至少其中一个是写入时


即使数据竞态在给定的编译器和平台下不存在问题, 也应该不能为存在数据竞态
找借口，根本没有所谓温和的数据竞态。
所有应该避免数据竞态, 如下方法:
- 不要修改变量
- 避免从多个 goroutine 访问同一个变量(管道，流水线访问，对变量的访问都是串行的，
先受限与流水线的一步，再受限于下一步，以此类推-- 串行受限)
- 允许多个 goroutine 访问同一个变量，但在同一时间只能有一个 goroutine
可以访问(互斥机制)


*/

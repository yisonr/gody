package main

/*
 * // go 语言实战， P83, go 语言的类型系统
 * go 语言是一种静态类型的编程语言，即编译器在编译时需知晓程序里每个值的类型，
 * 如果提前知道类型信息，编译器就可以确保程序合理使用值，这有助于减少潜在的
 * 内存异常和 bug, 并且使编译器对代码进行一些性能优化，提高执行效率
 *
 * 值的类型为编译器提供两部分信息:
 * - 值的规模(需要分配多少内存给这个值)
 * - 这段内存表示什么
 * 对于许多内置类型来说，规模和表示是类型名的一部分:
 * int64 类型的值需要8字节(64位), 表示一个整数
 * float32 类型的值需要4字节(32位), 表示一个IEEE-754定义的二进制浮点数
 * bool 类型的值需要1字节(8位), 表示布尔值 true 和 false
 *
 * 有些类型的内部表示与编译代码的机器的体系结构有关，如 int 值的大小根据
 * 编译所在的机器的体系结构，可能是8字节也可能是4字节; go 语言所有的引用类型，
 * 也与体系结构相关.
 *
 */

// 当用户声明一个新类型时(type typeName struct)，这个声明即为编译器提供了
// 一个框架，告知必要的内存大小和表示信息，声明后的类型与内置类型的运作方
// 式类似; 当声明一个变量时，这个变量对应的值总是会被初始化，这个值要么用
// 指定的值初始化，要么用零值(即变量类型的默认值) 做初始化
// 如果要初始化为某个非零值，就需要配合结构字面量和短变量声明操作符来创建变量

// 另一种声明用户定义的类型的方法是，基于一个已有的类型，将其作为新类型的类型
// 说明，标准库使用这种声明类型的方法，从内置类型创建出更多更加明确的类型，
// 并赋予更高级的功能, 如:
// type Duration int64
// int64 称为 Duration 的基础类型，但这两个类型是完全不同的有区别的类型
//

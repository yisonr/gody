package main

// 解码 S 表达式

/*
 * 对于标准库 encoding/... 提供的每一个 Marshal 函数，都有一个对应的 Unmarshal
 * 函数来解码
 * 如下:
  data :=  []byte("...")
  var movie Movie
  err := json.Unmarshal(data, &movie)

 * Unmarshal 函数使用反射修改已存在的 Movie 变量的字段，根据 Movie
 * 类型和输入数据来创建新的 map, 结构和 slice
 *
*/

// TODO: 读 json.Unmarshal 源码
// 为 S 表达式实现一个简单的 Unmarshal 函数，这个函数与 json.Unmarshal
// 函数类似，与 ./sexpr1.go 中的 Marshal 函数正好相反
// 注意: 一个鲁棒且通用的实现需要的代码量很大，所以此示例仅支持 S 表达式
// 一个有限的子集，并且没有优雅的处理错误, 示例的目的是阐释反射

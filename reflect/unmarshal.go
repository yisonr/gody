package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"text/scanner"
)

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

// Unmarshal 解析表达式数据并填充到非 nil 指针 out 指向的变量
func Unmarshal(data []byte, out interface{}) (err error) {
	lex := &lexer{
		scan: scanner.Scanner{Mode: scanner.GoTokens},
	}
	lex.scan.Init(bytes.NewReader(data))
	lex.next() // 获取第一个标记
	defer func() {
		if x := recover(); x != nil {
			err = fmt.Errorf("error at %s:%v", lex.scan.Position, x)
		}
	}()
	read(lex, reflect.ValueOf(out).Elem()) // Elem() 可寻址
	return nil
}

type Movie struct {
	Title, Subtitle string
	Year            int
	Color           bool
	Actor           map[string]string
	Oscars          []string
	Sequel          *string
}

func main() {
	filename := "./decodeText.json"
	fe, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		panic(err.Error())
	}
	defer fe.Close()
	byteValue, err := ioutil.ReadAll(fe)
	if err != nil {
		panic(err.Error())
	}
	var movie Movie
	err = Unmarshal(byteValue, &movie)
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println(movie)
	}
}

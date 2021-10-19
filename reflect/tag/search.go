package main

/*
 * 使用反射访问结构体字段标签
 *
 * 本代码定义一个工具函数 Unpack , 使用结构体字段标签来简化
 *  HTTP 处理程序
 * 首先调用 req.ParseForm() 来解析请求，在这之后 req.Form
 * 就有了所有的请求参数, 这个方法对 http GET 和 POST 请求都适用
 *
 * 示例了 reflect.Append, reflect.Value.Set, reflect.Value.Elem 等
 * 方法的使用, 以及 http.Request 解析请求，获取请求参数的方法
 * 包括读取结构体标签
 */

import (
	"fmt"
	"net/http"
	"reflect"
)

// search 用于处理 /search URL endpoint
func search(resp http.ResponseWriter, req *http.Request) {
	var data struct {
		Labels     []string `http:"l"`
		MaxResults int      `http:"max"`
		Exact      bool     `http:"x"`
	}
	data.MaxResults = 10 // 设置默认值
	if err := Unpack(req, &data); err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest) // 400
		return
	}
	// ...
	fmt.Fprintf(resp, "Search: %+v\n", data)
	err, query := Pack(reflect.ValueOf(data))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(query)
}

func main() {
	http.HandleFunc("/search", search)
	address := "127.0.0.1:8080"
	fmt.Printf("Start server at: %v\n", address)
	if err := http.ListenAndServe(address, nil); err != nil {
		fmt.Printf("Start server error: %v\n", err.Error())
	}
}

/* terminal1:
> go run pack.go params.go search.go
Start server at: 127.0.0.1:8080
map[l:<[]string Value> max:<int Value> x:<bool Value>]
&l=&l=s&l=p&max=10&x=true
*/

/* terminal2:
> curl "http://127.0.0.1:8080/search?l=s&l=p&max=10&x=true"
Search: {Labels:[s p] MaxResults:10 Exact:true}
*/

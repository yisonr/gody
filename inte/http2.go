package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	db := database{"shoes": 50, "socks": 5}
	mux := http.NewServeMux()
	// HandlerFunc 是一个类型, 这里是类型转换而非函数调用
	// type HandlerFunc func(ResponseWriter, *Request)
	// func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
	// 	f(w, r)
	// }
	// 它不仅是一个函数类型, 还拥有自己的方法，也满足接口 http.Hander
	// 它的 ServeHTTP 方法就是调用它本身
	// 所以 HandlerFunc 就是一个让函数值满足接口的一个适配器
	// 就是让函数 list 拥有了 ServeHTTP 方法, 实现了 Hander 接口
	mux.Handle("/list", http.HandlerFunc(db.list))
	mux.Handle("/price", http.HandlerFunc(db.price))
	// 由此可以用不同的方式满足 http.Hander 接口
	log.Fatal(http.ListenAndServe("localhost:8000", mux))
}

type dollars float32

func (d dollars) String() string { return fmt.Sprintf("$%.2f", d) }

type database map[string]dollars

func (db database) list(w http.ResponseWriter, req *http.Request) {
	for item, price := range db {
		fmt.Fprintf(w, "%s:%s\n", item, price)
	}
}

func (db database) price(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	price, ok := db[item]
	if !ok {
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "no such item:%q\n", item)
		return
	}
	fmt.Fprintf(w, "%s\n", price)
}

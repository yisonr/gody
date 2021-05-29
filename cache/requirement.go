package main

import (
	"io/ioutil"
	"net/http"
)

func httpGetBody(url string) (interface{}, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

//!-httpRequestBody

var HTTPGetBody = httpGetBody

func incomingURLs() <-chan string {
	ch := make(chan string)
	//-------------------------------------------
	// 函数返回一个只接收的通道(单向通道)
	// 同时注册一个 goroutine 往该通道发送数据
	// 函数的调用者只负责接收数据
	// 优点:
	// 通道在函数内维护，比较适合
	// 包级别的调用场景, 保证list的同步访问
	// 对于包内其他的复杂的处理任务，此种方式
	// 显得很妙
	go func() {
		for _, url := range []string{
			"https://baidu.com",
			"https://bing.com",
			"https://music.163.com",
			"https://baidu.com",
			"https://bing.com",
			"https://baidu.com",
			"https://music.163.com",
			"https://bing.com",
			"https://baidu.com",
			"https://bing.com",
		} {
			ch <- url
		}
		close(ch)
	}()
	//-------------------------------------------
	return ch
}

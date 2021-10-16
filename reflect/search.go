package main

import "net/http"

// search 用于处理 /search URL endpoint
func search(resp http.ResponseWriter, req *http.Request) {
	var data struct {
		Labels     []string `http:"1"`
		MaxResults int      `http:"max"`
		Exact      bool     `http:"x"`
	}
	data.MaxResults = 10 // 设置默认值
	if err := Unpack(req, &data); err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest) // 400
		return
	}

}

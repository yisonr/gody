package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"path"
)

func fetch(url string) (filename string, n int64, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	local := path.Base(resp.Request.URL.Path)
	if local == "/" || local == "." {
		local = "index.html"
	}
	f, err := os.Create(local)
	if err != nil {
		return "", 0, err
	}
	// 在这里试图使用延迟调用f.Close 关闭一个本地文件会存在问题,
	// 因为 os.Create 打开了一个文件对其进行写入, 创建;
	// 在许多文件系统中(TODO: linux文件系统), 尤其是 NFS中使用延迟写入(TODO),
	// 写错误往往不是立即返回而是推迟到文件关闭的时候, 如果无法检查关闭操作
	// 的结果, 就会导致一系列的数据丢失
	// 上述情况很少会遇到, 带熟悉linux文件系统以及 os.Create 实现的系统
	// 调用函数后再做研究

	n, err = io.Copy(f, resp.Body)
	// 关闭文件, 并保留错误消息
	if closeErr := f.Close(); err == nil {
		err = closeErr
	}
	return local, n, err
}

func main() {
	_, _, err := fetch(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
}

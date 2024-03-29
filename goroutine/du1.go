package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

func walkDir(dir string, fileSizes chan<- int64) {
	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			subdir := filepath.Join(dir, entry.Name())
			walkDir(subdir, fileSizes)
		} else {
			fileSizes <- entry.Size()
		}
	}
}

func dirents(dir string) []os.FileInfo {
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		// fmt.Fprintf(os.Stderr, "du1: %v\n", err)
		return nil
	}
	return entries
}

var verbose = flag.Bool("v", false, "show verbose progress messages")

func main() {
	// 确定初始目录
	flag.Parse()
	roots := flag.Args()
	// 默认文件所在目录
	if len(roots) == 0 {
		roots = []string{"."}
	}

	// 遍历文件树
	fileSizes := make(chan int64, 10)
	go func() {
		for _, root := range roots {
			walkDir(root, fileSizes) // 这里可以并发
		}
		close(fileSizes)
	}()

	// 定期输出结果
	var tick <-chan time.Time // 只接收通道
	if *verbose {
		// 如果命令行参数中 -v 标识没有指定，tick通道仍然
		// 是 nil, 它对应的情况在select中实际上被禁用
		tick = time.Tick(1 * time.Second)
	}
	var nfiles, nbytes int64
loop: // break 标签
	for {
		// 同时满足时的随机性，导致第一个 case 偶尔存在未执行的可能
		// 导致计算不准确???
		select {
		case size, ok := <-fileSizes:
			if !ok {
				break loop //  fileSizes关闭
				// 跳出select和for循环的逻辑
				// 没有标签的break只能跳出select的逻辑
			}
			nfiles++
			nbytes += size
		case <-tick:
			printDiskUsage(nfiles, nbytes)
		}
	} // 标签break退出for循环
	printDiskUsage(nfiles, nbytes)
}

func printDiskUsage(nfiles, nbytes int64) {
	// fmt.Println(nbytes)
	fmt.Printf("%d files %.1f GB\n", nfiles, float64(nbytes)/1e9)
	// TODO
	// 更好的显示文件大小
	// 根目录大小异常
}

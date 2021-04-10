package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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
		fmt.Fprintf(os.Stderr, "du1: %v\n", err)
		return nil
	}
	return entries
}

func main() {
	// 确定初始目录
	flag.Parse()
	roots := flag.Args()
	// 默认文件所在目录
	if len(roots) == 0 {
		roots = []string{"."}
	}

	// 遍历文件树
	fileSizes := make(chan int64)
	go func() {
		for _, root := range roots {
			walkDir(root, fileSizes)
		}
		close(fileSizes)
	}()

	// 输出结果
	var nfiles, nbytes int64
	for size := range fileSizes {
		nfiles++
		nbytes += size
	}
	printDiskUsage(nfiles, nbytes)

	// 定期输出结果
}

func printDiskUsage(nfiles, nbytes int64) {
	fmt.Println(nbytes)
	fmt.Printf("%d files %.1f GB\n", nfiles, float64(nbytes)/1e9)
	// TODO? 根文件夹的文件大小异常

}

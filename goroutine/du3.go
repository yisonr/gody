package main

// TODO: 加深理解
// 在du2的基础上添加取消功能
// 一个goroutine无法直接终止另一个，否则会让共享变量状态处于不确定状态???
// 通常任何时刻都很难知道有多少goroutine正在工作.
// 对于取消操作，需要一个可靠的机制在一个通道上广播一个事件，
// 这样很多goroutine可以检测到事件的发生
// 创建这样一个广播机制: 不在通道上发送值，而是关闭它
// 它的关闭表明程序需要停止它正在做的事

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"
)

func walkDir(dir string, n *sync.WaitGroup, fileSizes chan<- int64) {
	defer n.Done()
	// 1.检查是否取消
	if cancelled() {
		return
	}
	// 2.在 dirents 中检查是否取消
	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			n.Add(1) // 在每一个goroutine创建之前add(1), 然后在goroutine中
			// 传入sync.WaitGroup，使用 done() 在goroutine完成后减1
			subdir := filepath.Join(dir, entry.Name())
			go walkDir(subdir, n, fileSizes)
		} else {
			fileSizes <- entry.Size()
		}
	}
}

// sema 是一个用于限制目录并发数的计数信号量
var sema = make(chan struct{}, 20)

func dirents(dir string) []os.FileInfo {
	// sema <- struct{}{}        // 获取令牌  p操作
	// todo:
	// 性能剖析揭示了它的瓶颈在于此处获取信号量令牌的操作
	// 这里的select让取消操作的延迟从数百毫秒减为几十毫秒
	select {
	case sema <- struct{}{}: // 获取令牌 p操作
	case <-done: // 2. 检查是否取消
		return nil // 取消
	}
	defer func() { <-sema }() // 释放令牌 v操作, defer 后需是一个函数调用

	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		// fmt.Fprintf(os.Stderr, "du1: %v\n", err)
		return nil
	}
	return entries
}

var done = make(chan struct{})

func cancelled() bool {
	select {
	case <-done:
		return true
	default:
		return false
	}
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

	// 此goroutine读取标准输入, 它连接到终端,
	// 一旦开始读入输入(回车), 这个goroutine通过
	// 关闭done通道来广播事件.
	go func() {
		os.Stdin.Read(make([]byte, 1)) // 读一个字节
		close(done)
	}()

	// 遍历文件树
	fileSizes := make(chan int64, 10)
	var n sync.WaitGroup // 防goroutine泄露
	for _, root := range roots {
		n.Add(1)
		// 为每一个 walkDir 的调用创建一个新的 goroutine
		// 使用 sync.WaitGroup 为当前存活的 walkDir 调用计数
		go walkDir(root, &n, fileSizes)
	}
	go func() {
		n.Wait()
		// sync.WaitGroup 计数器减为0的时候，关闭fileSizes通道
		// 关闭通道是为了通知接收goroutine(在这里是主goroutine)，
		// 发送工作已经完毕
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
	// 主 goroutine 负责:
	// 1. 定时输出结果
	// 2. 检查来自标准输入的取消操作
	// 3. 检查 fileSizes 的关闭并退出程序, 否则继续从通道接收，并累积结果
loop: // break 标签, 这里的退出循环条件是: 1.回车退出 2.目录遍历完成
	for {
		// 同时满足时的随机性，导致第一个 case 偶尔存在未执行的可能
		// 导致计算不准确???
		select {
		case <-done:
			// 耗尽fileSizes以允许已有的goroutine结束, 否则产生泄露
			// 已存在的 goroutine 发送到通道的值需要被接收,
			// 否则发送 goroutine 会阻塞
			for range fileSizes {
				//不执行任何操作
			}
			return
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

	// 程序退出前总是输出
	printDiskUsage(nfiles, nbytes)
}

func printDiskUsage(nfiles, nbytes int64) {
	// fmt.Println(nbytes)
	fmt.Printf("%d files %.1f GB\n", nfiles, float64(nbytes)/1e9)
	// TODO
	// 更好的显示文件大小, Kb, Mb, Gb
	// 根目录大小异常
}

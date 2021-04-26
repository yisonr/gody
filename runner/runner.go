package main

import (
	"errors"
	"os"
	"os/signal"
	"time"
)

// runner 在给定的超时时间内执行一组任务，
// 并且在操作系统发送中断信号时结束这些任务
type Runner struct {
	// interrupt 通道报告从操作系统
	// 发送的信号
	interrupt chan os.Signal
	// complete 通道报告处理任务已经完成
	complete chan error
	// timeout 报告处理任务已经超时
	timeout <-chan time.Time
	// tasks 持有一组以索引顺序依次执行的函数
	tasks []func(int) // 参数为 int 的函数切片
}

// ErrTimeout 会在任务执行超时时返回
var ErrTimeout = errors.New("received timeout")

// ErrInterrupt 会在接收到操作系统的事件时返回
var ErrInterrupt = errors.New("received interrupt")

// 返回一个新的准备使用的 runner ,
// todo: 工厂函数
func New(d time.Duration) *Runner {
	return &Runner{
		interrupt: make(chan os.Signal, 1),
		complete:  make(chan error),
		timeout:   time.After(d),
		// tasks 不需要初始化，因为其是一个切片, 其零值是 nil
	}
}

// Add 将一个任务附加到 Runner 上，这个任务是一个
// 接收一个 int 类型的 ID 作为参数的函数
func (r *Runner) Add(tasks ...func(int)) {
	r.tasks = append(r.tasks, tasks...)
}

// Start 执行所有任务，并监视通道事件
func (r *Runner) Start() error {
	// 接收所有的中断信号
	signal.Notify(r.interrupt, os.Interrupt)
	// todo: Notify 函数

	// goroutine 执行任务
	go func() {
		r.complete <- r.run()
	}()

	select {
	case err := <-r.complete:
		return err // 任务顺利完成/出现系统中断
	case <-r.timeout: // 超时
		return ErrTimeout
	}
}

func (r *Runner) run() error {
	for id, task := range r.tasks {
		// 检测操作系统的中断信号
		if r.gotInterrupt() {
			return ErrInterrupt
		}
		// 执行已经注册的任务
		task(id)
	}
	return nil
}

// 验证是否接收到了中断信号
func (r *Runner) gotInterrupt() bool {
	select {
	case <-r.interrupt:
		// 停止接收之后的任何信号
		signal.Stop(r.interrupt)
		return true
	default:
		return false
	}
}

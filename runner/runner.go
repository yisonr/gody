package main

/* runner 包使用通道监视程序的执行时间，如果程序运行时间太长，也可以使用
* runner 包来终止程序, 当开发需要调度后台处理任务的程序的时候，这种模式
* 优势明显
* runner 包管理处理任务的运行和生命周期
 */

import (
	"errors"
	"os"
	"os/signal"
	"time"
)

// 程序展示了依据调度运行的无人值守的面向任务的程序，及其所使用的并发模式
// 在设计上，可支持以下终止点:
// - 程序在分配的时间内完成工作，正常终止
// - 程序没有及时完成工作，"自杀"
// - 接收到操作系统发送的中断事件，程序立刻试图清理状态并停止工作

// runner 在给定的超时时间内执行一组任务，
// 并且在操作系统发送中断信号时结束这些任务
type Runner struct {
	// interrupt 通道报告从操作系统
	// 发送的信号
	interrupt chan os.Signal
	/*
	* Signal 描述操作系统发送的信号，其底层实现通常会依赖操作系统的具体实现
	* 在 Unix 系统上是 syscall.Signal
	 */
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
		// 通道 interrupt 被初始化为缓冲区容量为1的通道，
		// 可以保证通道至少能接收一个来自语言运行时的
		// os.Signal 值，确保运行时发送这个事件的时候
		// 不会被阻塞
		complete: make(chan error),
		timeout:  time.After(d),
		// tasks 不需要初始化，因为其是一个切片, 其零值是 nil,
		// 满足初始化的要求，所以没有被明确初始化
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
	// Notify函数让signal包将输入信号转发到c。如果没有列出要传递的信号，会将所有输入信号传递到c；否则只传递列出的输入信号。
	// signal包不会为了向c发送信息而阻塞（就是说如果发送时c阻塞了，signal包会直接放弃
	// 调用者应该保证c有足够的缓存空间可以跟上期望的信号频率。对使用单一信号用于通知的通道，缓存为1就足够了。

	// goroutine 执行任务
	go func() {
		// 在此  goroutine 内发送
		r.complete <- r.run()
	}()

	select {
	case err := <-r.complete:
		// 在此  goroutine 内接收
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

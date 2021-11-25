package main

/*
 * runner 包展示如何使用通道来监视程序的执行时间, 当开发需要调度后台处理任务的
 * 程序时, 这种模式很有用. 此程序可能会作为 cron 作业执行, 或者在基于定时任务
 * 的云环境(如 iron.io)里执行
 */
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
		// 缓冲区容量为1, 可以保证通道至少能接收一个来自语言运行时
		// 的 os.Signal 值, 确保语言运行时发送这个事件的时候不会被阻塞
		complete: make(chan error),
		timeout:  time.After(d),
		// tasks 字段的零值是 nil, 满足初始化的要求, 所以不需要被明确初始化
	}
}

// Add 将一个任务附加到 Runner 上，这个任务是一个
// 接收一个 int 类型的 ID 作为参数的函数
func (r *Runner) Add(tasks ...func(int)) {
	r.tasks = append(r.tasks, tasks...)
}

// Start 执行所有任务，并监视通道事件
func (r *Runner) Start() error {
	// TODO: golang 信号处理
	// 接收所有的中断信号
	signal.Notify(r.interrupt, os.Interrupt)
	// golang中对信号的处理主要使用os/signal包中的两个方法:
	// - Notify 方法用来监听收到的信号
	// - Stop方法用来取消监听

	/*
	 *
	 * Notify 函数让 signal 包将输入信号转发到 channel, 如果没有列出要传递
	 * 的信号, 会将所有输入信号传递到 channel; 否则只传递列出的输入信号,
	 * signal 包不会为了向 channel 发送信息而阻塞, 即如果发送时 channel 阻塞了,
	 * signal 包会直接放弃): 调用者应该保证 channel 有足够的缓存空间可以跟上
	 * 期望的信号频率, 对使用单一信号用于通知的通道，缓存为1就足够了.
	 *
	 */

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
		// 联系 line 39 理解:
		// r.interrupt 缓冲区为1, 即在 select 执行到 default 分支后,
		// 恰好发生系统中断, 此时 r.interrupt 可以接收, 等下次 select 语句
		// 执行时, 就能进入 r.interrupt 分支
		// 如果缓冲区为0, 没做好接收准备(没执行到select)的情况下,
		// Notify 会放弃发送中断信号(line, 67), 导致中断信号遗漏

		// 停止接收之后的任何信号
		// TODO: 如果不停止接收的话, 会如何? 如何模拟信号的收发
		signal.Stop(r.interrupt)
		return true
	default:
		// default 分支会将接收 interrupt 的通道的阻塞调用转变为非阻塞的
		return false
	}
}

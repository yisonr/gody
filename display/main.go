package main

import (
	"os"
	"reflect"
)

type Movie struct {
	Title, Subtitle string
	Year            int
	Color           bool
	Actor           map[string]string
	Oscars          []string
	Sequel          *string
}

func main1() {
	strangelove := Movie{
		Title:    "dr.liu",
		Subtitle: "How i leard",
		Year:     1984,
		Color:    false,
		Actor: map[string]string{
			"Dr.Liu":       "peter sellrs",
			"Grp. Linoeel": "george select",
			"Bring. Gen":   "Slim Pickens",
		},
		Oscars: []string{
			"Best Actor",
			"Best Adapted Nomin",
			"Bsd shsdb",
		},
	}
	Display("strangelove", strangelove)

	Display("os.Stderr", os.Stderr)
	// Display os.Stderr (*os.File):
	// 输出:
	// (*(*os.Stderr).file).pfd.fdmu.state=0
	// (*(*os.Stderr).file).pfd.fdmu.rsema=0
	// (*(*os.Stderr).file).pfd.fdmu.wsema=0
	// (*(*os.Stderr).file).pfd.Sysfd=2
	// (*(*os.Stderr).file).pfd.pd.runtimeCtx=0
	// (*(*os.Stderr).file).pfd.iovecs=nil
	// (*(*os.Stderr).file).pfd.csema=0
	// (*(*os.Stderr).file).pfd.isBlocking=1
	// (*(*os.Stderr).file).pfd.IsStream=true
	// (*(*os.Stderr).file).pfd.ZeroReadIsEOF=true
	// (*(*os.Stderr).file).pfd.isFile=true
	// (*(*os.Stderr).file).name="/dev/stderr"
	// (*(*os.Stderr).file).dirinfo=nil
	// (*(*os.Stderr).file).nonblock=false
	// (*(*os.Stderr).file).stdoutOrErr=true
	// (*(*os.Stderr).file).appendMode=false
	// 即使非导出字段在反射下也是可见的

	// 也可把 Display 作用于 reflect.Value 上
	Display("rV", reflect.ValueOf(os.Stderr))
}

func main2() {
	var i interface{} = 3

	Display("i", i)
	// Display i (int):
	// i=3

	Display("&i", &i)
	// Display &i (*interface {}):
	// (*&i).type=int
	// (*&i).value=3
	/*
	 * TODO:
	 *  在 Display 函数的 Ptr 分支中，会调用这个值的 Elem 方法，返回一个代表
	 *  变量 i 的 Value, 其类型为 interface , 类似这种间接获得的 Value 可以
	 *  代表任何值, 包括接口
	 *  这时 display 函数递归调用自己, 输出接口的动态类型和动态值
	 */

}

func main() {
	/*
	 * 在以上 main1 和 main2 的实现中， Display 在对象图中存在循环引用时不会
	 * 自行终止, 如处理一个首尾相接的链表时:
	 *
	 */

	// 一个指向自己的结构体
	type Cycle struct {
		Value int
		Tail  *Cycle
	}

	var c Cycle
	c = Cycle{42, &c}
	Display("c", c) // Display 会输出一个持续增长的展开式
	// Display c (main.Cycle):
	// c.Value=42
	// (*c.Tail).Value=42
	// (*(*c.Tail).Tail).Value=42
	// (*(*(*c.Tail).Tail).Tail).Value=42
	// (*(*(*(*c.Tail).Tail).Tail).Tail).Value=42
	// (*(*(*(*(*c.Tail).Tail).Tail).Tail).Tail).Value=42
	// ........ 无穷无尽

	/*
	 * 很多 go 程序都至少包含一些循环引用的数据，让 Display 能鲁棒(健壮)的
	 * 处理这些循环引用需要一些技巧, 即需要记录曾经被访问过的引用, 当然成本
	 * 也不低, 一个通用的解决方案是使用 unsafe 语言特性(TODO: unsafe)
	 *
	 * TODO: 实验
	 * 循环引用在 fmt.Sprintf 中不构成大问题， 因为它很少尝试输出整个结构体，
	 * 比如，当输入是一个指针时, 它会输出指针的数字值, 这即打破了循环引用,
	 * 但如果遇到一个 slice 或 map 包含自身，还是会卡住，只是不值得为了这种
	 * 罕见的案例去承担处理循环引用的成本
	 *
	 */
}

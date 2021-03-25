package expr

import (
	"fmt"
	"math"
	"strings"
)

// Expr 算术表达式
type Expr interface {
	// Eval 返回表达式在 env 上下文的值
	Eval(Env) float64
	// Check 报告表达式中的错误，并把表达式中的
	// 变量加入到 vars 中
	Check(vars map[Var]bool) error
}

// Var 表示一个变量，比如 x
type Var string

func (v Var) Eval(env Env) float64 {
	return env[v]
}
func (v Var) Check(vars map[Var]bool) error {
	vars[v] = true
	return nil
}

// literal 是一个数字常量，比如 3.141
type literal float64

// _ 作为实参是为了满足接口定义的方法的形参
// func (l literal) Eval(Env) float64 {
// 	return float64(l)
// }
func (l literal) Eval(_ Env) float64 {
	return float64(l)
}

// 方法的接收者可以不表示，在不需要引用的时候
// func (_ literal) Check(vars map[Var]bool) error {
// 	return nil
// }
func (literal) Check(map[Var]bool) error {
	return nil
}

// unary 表示一元操作符表达式，比如 -x
type unary struct {
	op rune // '+', '-' 中的一个
	x  Expr
}

func (u unary) Eval(env Env) float64 {
	switch u.op {
	case '+':
		return +u.x.Eval(env)
	case '-':
		return -u.x.Eval(env)
	default:
		panic(fmt.Sprintf("unsupported unary operator: %q", u.op))
	}
}
func (u unary) Check(vars map[Var]bool) error {
	if !strings.ContainsRune("+-", u.op) {
		return fmt.Errorf("unexpected unary op %q", u.op)
	}
	return u.x.Check(vars)
}

// binary 表示二元操作符表达式， 比如 x+y
type binary struct {
	op   rune // '+', '-', '/' 中的一个
	x, y Expr
}

func (b binary) Eval(env Env) float64 {
	switch b.op {
	case '+':
		return b.x.Eval(env) + b.y.Eval(env)
	case '-':
		return b.x.Eval(env) - b.y.Eval(env)
	case '/':
		return b.x.Eval(env) / b.y.Eval(env)
	default:
		panic(fmt.Sprintf("unsupported binary operator: %q", b.op))
	}
}

func (b binary) Check(vars map[Var]bool) error {
	if !strings.ContainsRune("+-", b.op) {
		return fmt.Errorf("unexpected binary op %q", b.op)
	}
	if err := b.x.Check(vars); err != nil {
		return err
	}
	return b.y.Check(vars)
}

// call 表示函数调用表达式，比如 sin(x)
type call struct {
	fn   string
	args []Expr
}

func (c call) Eval(env Env) float64 {
	switch c.fn {
	case "pow":
		return math.Pow(c.args[0].Eval(env), c.args[1].Eval(env))
	case "sin":
		return math.Sin(c.args[0].Eval(env))
	case "sqrt":
		return math.Sqrt(c.args[0].Eval(env))
	default:
		panic(fmt.Sprintf("unsupported function call:%s", c.fn))
	}
}

var numParams = map[string]int{"pow": 2, "sin": 1, "sqrt": 1}

func (c call) Check(vars map[Var]bool) error {
	arity, ok := numParams[c.fn]
	if !ok {
		return fmt.Errorf("unknown function %q", c.fn)
	}
	if len(c.args) != arity {
		return fmt.Errorf("call to %s has %d args, want %d",
			c.fn, len(c.args), arity)
	}
	for _, arg := range c.args {
		if err := arg.Check(vars); err != nil {
			return err
		}
	}
	return nil
}

// 需要一个上下文环境把变量映射到数值
type Env map[Var]float64

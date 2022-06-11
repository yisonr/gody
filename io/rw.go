package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

// Writer is the interface that wraps the basic Write method.
//
// Write writes len(p) bytes from p to the underlying data stream.
// It returns the number of bytes written from p (0 <= n <= len(p))
// and any error encountered that caused the write to stop early.
// Write must return a non-nil error if it returns n < len(p).
// Write must not modify the slice data, even temporarily.
//
// Implementations must not retain p.
type Writer interface {
	Write(p []byte) (n int, err error)
	io.Reader
}

// Reader is the interface that wraps the basic Read method.
//
// Read reads up to len(p) bytes into p. It returns the number of bytes
// read (0 <= n <= len(p)) and any error encountered. Even if Read
// returns n < len(p), it may use all of p as scratch space during the call.
// If some data is available but not len(p) bytes, Read conventionally
// returns what is available instead of waiting for more.
//
// When Read encounters an error or end-of-file condition after
// successfully reading n > 0 bytes, it returns the number of
// bytes read. It may return the (non-nil) error from the same call
// or return the error (and n == 0) from a subsequent call.
// An instance of this general case is that a Reader returning
// a non-zero number of bytes at the end of the input stream may
// return either err == EOF or err == nil. The next Read should
// return 0, EOF.
//
// Callers should always process the n > 0 bytes returned before
// considering the error err. Doing so correctly handles I/O errors
// that happen after reading some bytes and also both of the
// allowed EOF behaviors.
//
// Implementations of Read are discouraged from returning a
// zero byte count with a nil error, except when len(p) == 0.
// Callers should treat a return of 0 and nil as indicating that
// nothing happened; in particular it does not indicate EOF.
//
// Implementations must not retain p.
type Reader interface { // TODO: 完全理解实现 Reader 接口的原则
	Read(p []byte) (n int, err error)
}

// 以下示例使用 bytes, fmt 和 os 包来进行缓冲, 拼接和写字符串到 stdout
func main() {
	// 创建一个 Buffer 值, 并将一个字符串写入 Buffer
	var b bytes.Buffer
	b.Write([]byte("Hello "))

	// 使用 Fprintf 将一个字符串拼接到 Buffer 里
	fmt.Fprintf(&b, "World!\n")

	// 将 Buffer 的内容输出到标准输出设备
	b.WriteTo(os.Stdout)
}

/*
	*bytes.Buffer 实现了 io.Writer 接口
	// write appends the contents of p to the buffer, growing the buffer as
	// needed. the return value n is the length of p; err is always nil. if the
	// buffer becomes too large, write will panic with errtoolarge.
	func (b *buffer) write(p []byte) (n int, err error) {
		b.lastread = opinvalid
		m, ok := b.trygrowbyreslice(len(p))
		if !ok {
			m = b.grow(len(p))
		}
		return copy(b.buf[m:], p), nil
	}

	TODO: 缓冲区扩容原理, 实验验证缓冲区太大, Write 崩溃的情况


	*os.File 实现了 io.Writer 接口
	// Write writes len(b) bytes from b to the File.
	// It returns the number of bytes written and an error, if any.
	// Write returns a non-nil error when n != len(b).
	func (f *File) Write(b []byte) (n int, err error) {
		if err := f.checkValid("write"); err != nil {
			return 0, err
		}
		n, e := f.write(b)
		if n < 0 {
			n = 0
		}
		if n != len(b) {
			err = io.ErrShortWrite
		}

		epipecheck(f, e)

		if e != nil {
			err = f.wrapErr("write", e)
		}

		return n, err
	}

*/

package main

// 编码 S 表达式

import (
	"bytes"
	"fmt"
	"reflect"
)

// 结构被编码为一个关于字段绑定(field binding) 的列表，每个字段绑定都是一个
// 两个元素的列表，其中第一个元素(使用符号) 是字段名，第二个元素是字段值，
// map 也编码为元素对的列表，每个元素对都是map中一项的键和值。
// S 表达式使用形式为(key . value) 的单个构造单元(cons cell) 来表示键值对,
// 而不是用双元素的列表?

func encode(buf *bytes.Buffer, v reflect.Value) error {
	switch v.Kind() {
	case reflect.Invalid:
		buf.WriteString("nil")

	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		fmt.Fprintf(buf, "%d", v.Int())

	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64:
		fmt.Fprintf(buf, "%d", v.Uint())

	case reflect.String:
		fmt.Fprintf(buf, "%q", v.String())

	case reflect.Ptr:
		// 递归指针的值
		return encode(buf, v.Elem())

	case reflect.Array, reflect.Slice:
		buf.WriteByte('(')
		for i := 0; i < v.Len(); i++ {
			if i > 0 {
				buf.WriteByte(' ')
			}
			if err := encode(buf, v.Index(i)); err != nil {
				return err
			}
		}
		buf.WriteByte(')')

	case reflect.Struct:
		buf.WriteByte('(')
		for i := 0; i < v.NumField(); i++ {
			if i > 0 {
				buf.WriteByte(' ')
			}
			fmt.Fprintf(buf, "(%s ", v.Type().Field(i).Name)
			if err := encode(buf, v.Field(i)); err != nil {
				return err
			}
			buf.WriteByte(')')
		}
		buf.WriteByte(')')

	case reflect.Map:
		buf.WriteByte('(')
		for i, key := range v.MapKeys() {
			if i > 0 {
				buf.WriteByte(' ')
			}
			buf.WriteByte('(')
			if err := encode(buf, key); err != nil {
				return err
			}
			buf.WriteByte(' ')
			if err := encode(buf, v.MapIndex(key)); err != nil {
				return err
			}
			buf.WriteByte(')')
		}
		buf.WriteByte(')')
	default: // float, complex, bool, chan, func, interface
		return fmt.Errorf("unsupported type:%s", v.Type())
	}
	return nil
}

// Marshal 把 go 的值编码为 S 表达式的形式
func Marshal(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if err := encode(&buf, reflect.ValueOf(v)); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

type Movie struct {
	Title, Subtitle string
	Year            int
	// Color           bool
	Actor  map[string]string
	Oscars []string
	Sequel *string
}

func main() {
	strangelove := Movie{
		Title:    "dr.liu",
		Subtitle: "How i leard",
		Year:     1984,
		// Color:    false,
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
	out, err := Marshal(strangelove)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("编码---------------------:\n")
	fmt.Println(string(out))
}

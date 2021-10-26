package main

import (
	"fmt"
	"reflect"
	"strings"
)

type ReqStruct struct {
	Labels     []string `http:"l"`
	MaxResults int      `http:"max"`
	Exact      bool     `http:"x"`
}

func Pack(v reflect.Value) (error, string) {
	var query string
	var err error
	switch v.Kind() {
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			// 可使用 reflect.StructField 获取关于 struct 的信息，包括 tag
			fieldInfo := v.Type().Field(i) // a reflect.StructField
			tag := fieldInfo.Tag           // a reflect.StructTag
			fieldKey := tag.Get("http")
			if fieldKey == "" {
				fieldKey = strings.ToLower(fieldInfo.Name)
			}
			value := v.Field(i)
			fieldKey = "&" + fieldKey
			err, query = decode(fieldKey, query, value)
			if err != nil {
				return err, query
			}
		}
		return nil, query

	default:
		return nil, fmt.Sprintf("Unsupported type: %s", v.Type())
	}
}

func decode(fieldKey, name string, v reflect.Value) (error, string) {
	// TODO: 使用buf 代替 name
	name += fieldKey + "="
	switch v.Kind() {
	case reflect.Slice:
		var err error
		for i := 0; i < v.Len(); i++ {
			// err, name := decode(name, v.Index(i))
			// Err: 这种写法 err, name 都是 for {} 的局部变量
			// 即 line53 返回的值是 line40 的 name 值
			err, name = decode(fieldKey, name, v.Index(i))
			if err != nil {
				return err, name
			}
		}
		return nil, name
	case reflect.Int:
		name += fmt.Sprintf("%d", v.Int())
		return nil, name
	case reflect.Bool:
		name += fmt.Sprintf("%v", v.Bool())
		return nil, name
	case reflect.String:
		name += v.String()
		return nil, name

	default: //... 更多的结构
		return nil, fmt.Sprintf("Unsupported type: %s", v.Type())
	}
}

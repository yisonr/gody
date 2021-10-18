package main

import (
	"fmt"
	"reflect"
)

type ReqStruct struct {
	Labels     []string `http:"l"`
	MaxResults int      `http:"max"`
	Exact      bool     `http:"x"`
}

func Pack(v reflect.Value) (error, string) {
	switch v.Kind() {
	case reflect.Struct:
		return nil, ""

	default:
		return nil, fmt.Sprintf("Unsupported type: %s", v.Type())
	}
}

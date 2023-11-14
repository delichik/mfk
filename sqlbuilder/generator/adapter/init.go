package adapter

import (
	"fmt"
	"reflect"
)

var typeStringer = reflect.TypeOf((*fmt.Stringer)(nil)).Elem()

var nilCheckHandlers = map[string]Handler{}
var toStringHandlers = map[string]Handler{}

type Handler func(fieldType reflect.Type) string

func NilCheck(fieldType reflect.Type) string {
	h, ok := nilCheckHandlers[NameOf(fieldType)]
	if ok {
		return h(fieldType)
	}
	return ""
}

func AsString(fieldType reflect.Type) string {
	h, ok := toStringHandlers[NameOf(fieldType)]
	if ok {
		return h(fieldType)
	} else if fieldType.Implements(typeStringer) {
		return "\"\\\"\" + ${fieldName}.String() + \"\\\"\""
	}
	return "fmt.Sprintf(\"%v\", ${fieldName})"
}

func NameOf(obj any) string {
	if t, ok := obj.(reflect.Type); ok {
		return t.PkgPath() + "/" + t.Name()
	}
	t := reflect.TypeOf(obj)
	return t.PkgPath() + "/" + t.Name()
}

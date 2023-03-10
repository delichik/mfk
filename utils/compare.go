package utils

import "reflect"

func DeepCompare[T comparable](a T, b interface{}) bool {
	tb, ok := b.(T)
	if !ok {
		panic("struct is not matched")
	}
	return reflect.DeepEqual(tb, a)
}

package adapter

import "reflect"

func init() {
	nilCheckHandlers[reflect.Pointer.String()] = nilCheck_Pointer
	toStringHandlers[reflect.Pointer.String()] = toString_Pointer

	nilCheckHandlers[NameOf(1)] = nilCheck_Int
	toStringHandlers[NameOf(1)] = toString_Int

	nilCheckHandlers[NameOf("")] = nilCheck_String
	toStringHandlers[NameOf("")] = toString_String
}

func nilCheck_Pointer(_ reflect.Type) string {
	return "${fieldName} != nil"
}

func toString_Pointer(value reflect.Type) string {
	value = value.Elem()
	return AsString(value)
}

func nilCheck_Int(_ reflect.Type) string {
	return "${fieldName} != 0"
}

func toString_Int(_ reflect.Type) string {
	return "strconv.Itoa(${fieldName})"
}

func nilCheck_String(_ reflect.Type) string {
	return "${fieldName} != \"\""
}

func toString_String(_ reflect.Type) string {
	return "\"\\\"\" + ${fieldName} + \"\\\"\""
}

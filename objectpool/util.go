package objectpool

import (
	"unsafe"
)

func cleanObject[T any](obj *T) {
	lowData := (*[unsafe.Sizeof(*obj)]byte)(unsafe.Pointer(obj))
	for i := range *lowData {
		(*lowData)[i] = byte(0)
	}
}

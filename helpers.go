package goof

import (
	"unsafe"
)

func makeInterface(typ, val unsafe.Pointer) interface{} {
	return *(*interface{})(unsafe.Pointer(&[2]unsafe.Pointer{typ, val}))
}

func dataPtr(val interface{}) unsafe.Pointer {
	return (*[2]unsafe.Pointer)(unsafe.Pointer(&val))[1]
}

package goof

import (
	"reflect"
	"unsafe"
)

func makeInterface(typ, val unsafe.Pointer) interface{} {
	return *(*interface{})(unsafe.Pointer(&[2]unsafe.Pointer{typ, val}))
}

func dataPtr(val interface{}) unsafe.Pointer {
	return unsafe.Pointer(reflect.ValueOf(&val).Elem().InterfaceData()[1])
}

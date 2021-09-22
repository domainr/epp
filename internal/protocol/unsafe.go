package protocol

import (
	"reflect"
	"unsafe"
)

func unsafeBytes(s string) []byte {
	var b []byte
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh.Data = sh.Data
	bh.Len = sh.Len
	bh.Cap = bh.Len
	return b
}

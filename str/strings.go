package str

import (
	"reflect"
	"unsafe"
)

func Or(strs ...string)string{
	for _, v := range strs {
		if v != ""{
			return v
		}
	}
	return ""
}

func In(str string, strs ... string )bool{
	for _, v := range strs {
		if str == v{
			return true
		}
	}
	return false
}

func BytesOf(s string)[]byte{
	h:=(*reflect.SliceHeader)(unsafe.Pointer(&s))
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: h.Data,
		Len:  h.Len,
		Cap:  h.Len,
	}))
}

func StringOf(b []byte)string{
	return *(*string)(unsafe.Pointer(&b))
}

func And(ss ...string)bool{
	for _, s := range ss {
		if s ==""{
			return false
		}
	}
	return true
}


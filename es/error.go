package es

import (
	"fmt"
	"runtime"
)


func Errorf(f string, args ...interface{}) error {
	err := fmt.Errorf(f, args...)
	pc, _, _, _ := runtime.Caller(1)
	fun := runtime.FuncForPC(pc)
	if fun == nil{
		return fmt.Errorf(f,args...)
	}
	return fmt.Errorf("%s:%w", fun.Name(), err)
}

func WrapError( str string,err error) error {
	pc, file, line, _ := runtime.Caller(1)
	fun := runtime.FuncForPC(pc)
	if fun == nil{
		return fmt.Errorf("%s:%d %s: %w",file,line,str,err)
	}
	return fmt.Errorf("%s:%d %s: %s:\n%w",file,line, fun.Name(), str, err)
}


func AddCaller(err error)error{
	pc, file, line, _ := runtime.Caller(1)
	fun := runtime.FuncForPC(pc)
	if fun == nil{
		return fmt.Errorf("%s:%d %w",file,line,err)
	}
	return fmt.Errorf("%s:%d %s: \n%w",file,line, fun.Name(), err)
}
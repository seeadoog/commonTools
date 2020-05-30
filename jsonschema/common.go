package jsonschema

import (
	"fmt"
	"strings"
)

type Error struct {
	Path string
	Info string
}

type Validator interface {
	Validate(path string,value interface{},errs *[]Error)
}

func appendString(s ...string)string{
	sb:=strings.Builder{}
	for _, str:= range s {
		sb.WriteString(str)
	}
	return sb.String()
}

func panicf(f string,args ...interface{}){
	panic(fmt.Sprintf(f,args...))
}
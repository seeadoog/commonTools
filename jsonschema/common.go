package jsonschema

import (
	"fmt"
	"strconv"
	"strings"
)

type Error struct {
	Path string
	Info string
}

type Validator interface {
	Validate(path *pathTree,value interface{},errs *[]Error)
}

func appendString(s ...string)string{
	sb:=strings.Builder{}
	n:=0
	for i:=0;i< len(s);i++{
		n+= len(s[i])
	}
	sb.Grow(n)
	for i:=0;i< len(s);i++{
		sb.WriteString(s[i])
	}
	return sb.String()
}

func panicf(f string,args ...interface{}){
	panic(fmt.Sprintf(f,args...))
}

func String(v interface{})string{
	switch v.(type) {
	case string:
		return v.(string)
	}
	return fmt.Sprintf("%v",v)
}

func Number(v interface{})float64{
	switch v.(type) {
	case float64:
		return v.(float64)
	case bool:
		if v.(bool){
			return 1
		}
		return 0
	case string:
		i,_:=strconv.ParseFloat(v.(string),64)
		return i
	}
	return 0
}


func Equal(a,b interface{})bool{
	return String(a) ==String(b)
}
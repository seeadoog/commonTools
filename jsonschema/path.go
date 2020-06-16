package jsonschema

import "strings"

type path struct {
	parent *path
	path string
}
func (path *path)String()string{
	strs:=[]string{}
	p:=path
	for p != nil{
		strs = append(strs,p.path)
		p = p.parent
	}
}
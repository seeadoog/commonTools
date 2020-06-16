package jsonschema

import (
	"fmt"
	"strings"
	"testing"
)

func Test_path_Set(t *testing.T) {
	p:=&pathTree{
		path: "$",
	}

	r:=p
	r = r.AddChild("b")
	   fmt.Println(r.AddChild("1").String())
	   fmt.Println(r.AddChild("2").String())


}

type pathTree struct {
	path string
	parent *pathTree
	children []*pathTree
	ci int
}

func (p *pathTree)reset(){
	p.ci = 0
	p.parent = nil
}

func newPathTree()*pathTree{
	return new(pathTree)
}

func (p *pathTree)AddChild(path string)*pathTree{

	if p.ci >= len(p.children){
		old:=p.children
		length:= len(old)
		if length == 0{
			length = 8
		}
		p.children = make([]*pathTree, length*2)
		copy(p.children,old)
	}
	n:=newPathTree()
	n.parent = p
	n.path = path
	p.children[p.ci] = new(pathTree)
	p.ci++
	return n
}

func (p *pathTree)String()string{
	strs:=[]string{}
	r:=p
	for r != nil{
		strs = append(strs,r.path)
		r = r.parent
	}


	switch len(strs) {
	case 0:
		return ""
	case 1:
		return strs[0]

	}
	n:=len(strs)-1
	for i:=0 ;i< len(strs);i++{
		n+= len(strs[i])
	}
	bf:=strings.Builder{}
	bf.Grow(n)
	for i:= len(strs)-1;i>0;i--{
		bf.WriteString(strs[i])
		bf.WriteString(".")
	}
	bf.WriteString(strs[0])
	return bf.String()
}

package jsonschema

import (
	"bytes"
	"strings"
	"sync"
)

type path struct {
	keys []string
	depth int
}

func newPath()*path{
	return &path{
		keys: make([]string,4),

	}
}

func (p *path)Set(path string){
	if p.depth >= len(p.keys){
		old:=p.keys
		p.keys = make([]string, len(old)*2)
		copy(p.keys,old)
	}
	p.keys[p.depth] = path
}



func (p *path)String()string{
	bf:=bytes.Buffer{}
	for i:=0;i<=p.depth;i++{
		 bf.WriteString(p.keys[i])
		 bf.WriteString(".")
	}
	return bf.String()
}

func (p *path)Next(path string){
	p.Set(path)
	p.depth ++
}

func (p *path)Back(){
	p.depth --
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
var pool sync.Pool
func newPathTree()*pathTree{
	 t:=pool.Get()
	 if t == nil{
	 	return new(pathTree)
	 }
	 p:=t.(*pathTree)
	 p.reset()
	 return p
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
	if p.children[p.ci]!=nil{
		n:=p.children[p.ci]
		n.path = path
		n.parent = p
		p.ci++
		return n

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

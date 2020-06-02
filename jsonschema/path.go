package jsonschema

import "bytes"

type path struct {
	keys []string
	depth int
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
	for i:=0;i<p.depth;i++{
		 bf.WriteString(p.keys[i])
		 bf.WriteString(".")
	}
	if len(p.keys)>p.depth{
		bf.WriteString(p.keys[p.depth])
	}
	return bf.String()
}

func (p *path)Next(){
	p.depth ++
}

func (p *path)Back(){
	p.depth --
}


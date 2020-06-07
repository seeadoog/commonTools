package jsonschema

import "bytes"

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


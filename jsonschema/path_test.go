package jsonschema

import (
	"fmt"
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

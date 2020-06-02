package jsonschema

import (
	"fmt"
	"testing"
)

func Test_path_Set(t *testing.T) {
	p:=path{
		keys:  make([]string,4),
		depth: 0,
	}
	p.Set("$")
	p.Next()
		p.Set("1")
		p.Next()
			p.Set("2")
	fmt.Println(p.String())
		p.Back()
		p.Set("2")
	fmt.Println(p.String())
}
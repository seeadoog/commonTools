package ngcfg

import (
	"fmt"
	"testing"
)

func TestNewLinkedMap(t *testing.T) {
	m:=NewLinkedMap()

	m.Set("1",1)
	m.Set("2",1)
	m.Set("3",1)
	m.Set("4",1)

	e:=m.MapItem()
	for e!= nil{
		fmt.Println(e.Val,e.Key)
		e = e.Next()
	}
}
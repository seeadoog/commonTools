package lrucache

import (
	"fmt"
	"testing"
)

func Test_lrucacheImp_Get(t *testing.T) {
	l:=NewLRUCache(3)
	l.Set("1",1)
	l.Set("2",2)
	l.Set("3",3)
	fmt.Println(l.Get("1"))
	fmt.Println(l.Get("2"))
	fmt.Println(l.Get("3"))
	l.Set("4",4)
	fmt.Println(l.Get("4"))
	fmt.Println(l.Get("1"))
	fmt.Println(l.Get("2"))
	l.Set("5",5)
	fmt.Println(l.Get("3"))
}
package main

import (
	"fmt"
	"sync"
	"testing"
)

func TestMap(t *testing.T) {
	a := &R{}
	for i := 0; i < 10000000; i++ {
		a = &R{
			A: map[string]interface{}{},
			B: map[string]interface{}{},
			C: map[string]interface{}{},
		}
		//json.Unmarshal(js,a)

	}
	fmt.Println(a)
}

func clearMap(m map[string]interface{}) {
	for k := range m {
		delete(m, k)
	}
}

type R struct {
	A map[string]interface{} `json:"a"`
	B map[string]interface{} `json:"b"`
	C map[string]interface{} `json:"c"`
}

var js = []byte(`
{
  "a": {
    "1": 1,
    "2": "2"
  },
  "b": {
    "1": 1,
    "2": "2",
    "3": 1,
    "4": "2",
    "5": "5"
  },
  "c": {
    "2": "2",
    "3": 1,
    "4": "2"
  }
}
`)

func TestMap2(t *testing.T) {

	var pool = sync.Pool{}

	a := &R{}
	for i := 0; i < 1000000; i++ {
		m := pool.Get()
		if m == nil {
			a = &R{
				A: map[string]interface{}{
					"e": 1,
				},
				B: map[string]interface{}{},
				C: map[string]interface{}{},
			}
		} else {
			a = m.(*R)
			clearMap(a.A)
			clearMap(a.B)
			clearMap(a.C)
		}
		//fmt.Println(a.A)
		//json.Unmarshal(js,a)
		//fmt.Println(a.A)
		pool.Put(a)
	}
}

//go:generate stringer -type=S
type S []int

func TestStringer(t *testing.T) {

}

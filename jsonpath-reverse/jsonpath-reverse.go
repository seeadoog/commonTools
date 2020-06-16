package jsonpath_reverse

import (
	"fmt"
	"reflect"
)

type tokenItem struct {
	index int
	key   string
}

//a.b[5]
//a b [5]
type CompiledJsonPath struct {
	tkns []tokenItem
}

func (c *CompiledJsonPath) Get(in interface{}) (interface{}, error) {
	current := in
foreach:
	for _, tkn := range c.tkns {
		if tkn.index < 0 { // get in object
			switch current.(type) {
			case map[string]interface{}:
				current = current.(map[string]interface{})[tkn.key]
				continue foreach
			default:
				val := reflect.ValueOf(current)
				if val.Kind() == reflect.Map {
					cv := val.MapIndex(reflect.ValueOf(tkn.key))
					if cv.IsZero() {
						current = nil
						continue foreach
					}
					if cv.CanInterface() {
						current = cv.Interface()
						continue foreach
					}
				}
				return nil, fmt.Errorf("get object in not object value")
			}
		} else {
			switch current.(type) {
			case []interface{}:
				arr := current.([]interface{})
				if len(arr) > tkn.index {
					current = arr[tkn.index]
					continue
				}
				return nil, fmt.Errorf("index out of range")
			default:
				val := reflect.ValueOf(current)
				if val.Kind() == reflect.Slice {
					if val.Len() > tkn.index {
						cv:=val.Index(tkn.index)
						if cv.IsZero(){
							current = nil
							continue foreach
						}
						if cv.CanInterface(){
							current = cv.Interface()
							continue foreach
						}
					}
					return nil,fmt.Errorf("index out of range")
				}
				return nil,fmt.Errorf("cannot index at not array value")
			}

		}
	}
	return current,nil
}

func (c *CompiledJsonPath)Set(obj interface{},value interface{})error{
	switch obj.(type) {
	case map[string]interface{}:

	}
}

// c[0]
func (c *CompiledJsonPath)setValue(in interface{},value interface{})error{
	cu:= in
	var parent interface{}
	foreach:
	for i, tkn := range c.tkns {
		if i< len(c.tkns) -1{  // build mid path

			if tkn.index <0{
				switch cu.(type) {
				case map[string]interface{}:
					cum:=cu.(map[string]interface{})
					cu = cum[tkn.key]
					if cu == nil{
						cu = map[string]interface{}{}
						cum[tkn.key] = cu
					}
					continue foreach
				}
			}

		}
	}
}
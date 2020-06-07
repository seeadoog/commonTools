package jsonschema

import (
	"fmt"
	"reflect"
	"strconv"
)

type Type string

func (t Type)Validate(path *pathTree,value interface{},errs *[]Error){
	if value == nil{
		return
	}
	switch t {
	case "object":
		if _,ok:=value.(map[string]interface{});!ok{
			rt:=reflect.ValueOf(value)
			if rt.Kind() == reflect.Struct{
				return
			}
			if rt.Kind() == reflect.Ptr{
				if !rt.IsNil(){
					if rt.Elem().Kind() == reflect.Struct{
						return
					}
				}
			}
			*errs = append(*errs,Error{
				Path: path.String(),
				Info: "type must be object",
			})
		}
	case "string":
		if _,ok:=value.(string);!ok{
			*errs = append(*errs,Error{
				Path: path.String(),
				Info: "type must be string",
			})
		}
	case "number","integer":
		if _,ok:=value.(float64);!ok{
			rt:=reflect.TypeOf(value)
			switch rt.Kind() {
			case reflect.Int,reflect.Int16,reflect.Int8,reflect.Int32,reflect.Int64,reflect.Uint8,reflect.Uint16,reflect.Uint32,reflect.Uint64,reflect.Uint,reflect.Float32,reflect.Float64:
				return
			}
			*errs = append(*errs,Error{
				Path: path.String(),
				Info: "type must be number",
			})
		}
	case "boolean","bool":
		if _,ok:=value.(bool);!ok{
			*errs = append(*errs,Error{
				Path: path.String(),
				Info: "type must be boolean",
			})
		}
	case "array":
		if _,ok:=value.([]interface{});!ok{
			*errs = append(*errs,Error{
				Path: path.String(),
				Info: "type must be array",
			})
		}

	}
}

type MaxLength int

func (l MaxLength)Validate(path *pathTree,value interface{},errs *[]Error){

	switch value.(type) {
	case string:
		if len(value.(string)) > int(l){
			*errs = append(*errs,Error{
				Path: path.String(),
				Info: "length must be <= "+strconv.Itoa(int(l)),
			})
		}
	case []interface{}:
		if len(value.([]interface{})) >int(l){
			*errs = append(*errs,Error{
				Path: path.String(),
				Info: "length must be <= "+strconv.Itoa(int(l)),
			})
		}
	}

}

type MinLength int

func (l MinLength)Validate(path *pathTree,value interface{},errs *[]Error){
	switch value.(type) {
	case string:
		if len(value.(string)) < int(l){
			*errs = append(*errs,Error{
				Path: path.String(),
				Info: "length must be >= "+strconv.Itoa(int(l)),
			})
		}
	case []interface{}:
		if len(value.([]interface{})) <int(l){
			*errs = append(*errs,Error{
				Path: path.String(),
				Info: "length must be >= "+strconv.Itoa(int(l)),
			})
		}
	}
}

type Maximum float64

func (m Maximum) Validate(path *pathTree, value interface{}, errs *[]Error) {
	val,ok:=value.(float64)
	if !ok{
		return
	}
	if val >float64(m){
		*errs = append(*errs,Error{
			Path: path.String(),
			Info: appendString("value must be <=",strconv.FormatFloat(float64(m),'f',-1,64)),
		})
	}
}

type Minimum float64

func (m Minimum) Validate(path *pathTree, value interface{}, errs *[]Error) {
	val,ok:=value.(float64)
	if !ok{
		return
	}
	if val<float64(m){
		*errs = append(*errs,Error{
			Path: path.String(),
			Info: appendString("value must be >=",strconv.FormatFloat(float64(m),'f',-1,64)),
		})
	}
}


type Enums []interface{}

func (enums Enums) Validate(path *pathTree, value interface{}, errs *[]Error) {
	if value == nil{
		return
	}
	for _, e := range enums {
		if e == value{
			return
		}
	}

	for _,e:= range enums{
		if Equal(e,value){
			return
		}
	}
	*errs = append(*errs,Error{
		Path: path.String(),
		Info: fmt.Sprintf("must be one of %v",enums),
	})
}

type Required []string

func (r Required) Validate(path *pathTree, value interface{}, errs *[]Error) {
	m,ok:=value.(map[string]interface{})
	if !ok{
		return
	}
	for _,key := range r {
		if _,ok:=m[key];!ok{
			*errs = append(*errs,Error{
				Path: appendString(path.String(),".",key),
				Info: "field is required",
			})
		}
	}
}

//// 限定数组的长度
//type Length int
//
//func (l Length) Validate(path *pathTree, value interface{}, errs *[]Error) {
//	arr,ok:=value.([]interface{})
//	if !ok{
//		return
//	}
//	if len(arr)>int(l){
//		*errs = append(*errs,Error{
//			Path: path,
//			Info: appendString("length must be length than ",strconv.Itoa(int(l))),
//		})
//	}
//}
//
//func NewLength(i interface{},parent Validator)(Validator,error){
//	it,ok:=i.(float64)
//	if !ok{
//		return nil,fmt.Errorf("value of 'length' must be integer:%v",i)
//	}
//	return Length(it),nil
//}

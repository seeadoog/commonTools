package jsonschema

import (
	"fmt"
	"strconv"
)

type Type string

func (t Type)Validate(path string,value interface{},errs *[]Error){
	if value == nil{
		return
	}
	switch t {
	case "object":
		if _,ok:=value.(map[string]interface{});!ok{
			*errs = append(*errs,Error{
				Path: path,
				Info: "type must be object",
			})
		}
	case "string":
		if _,ok:=value.(string);!ok{
			*errs = append(*errs,Error{
				Path: path,
				Info: "type must be string",
			})
		}
	case "boolean","bool":
		if _,ok:=value.(bool);!ok{
			*errs = append(*errs,Error{
				Path: path,
				Info: "type must be boolean",
			})
		}
	case "number","integer":
		if _,ok:=value.(float64);!ok{
			*errs = append(*errs,Error{
				Path: path,
				Info: "type must be number",
			})
		}
	case "array":
		if _,ok:=value.([]interface{});!ok{
			*errs = append(*errs,Error{
				Path: path,
				Info: "type must be array",
			})
		}

	}
}

type MaxLength int

func (l MaxLength)Validate(path string,value interface{},errs *[]Error){
	str,ok:=value.(string)
	if !ok{
		return
	}
	if len(str) > int(l){
		*errs = append(*errs,Error{
			Path: path,
			Info: "length must be <= "+strconv.Itoa(int(l)),
		})
	}
}

type MinLength int

func (l MinLength)Validate(path string,value interface{},errs *[]Error){
	str,ok:=value.(string)
	if !ok{
		return
	}
	if len(str) < int(l){
		*errs = append(*errs,Error{
			Path: path,
			Info: "length must be > "+strconv.Itoa(int(l)),
		})
	}
}

type Maximum float64

func (m Maximum) Validate(path string, value interface{}, errs *[]Error) {
	val,ok:=value.(float64)
	if !ok{
		return
	}
	if val >float64(m){
		*errs = append(*errs,Error{
			Path: path,
			Info: appendString("value must be <=",strconv.FormatFloat(float64(m),'f',-1,64)),
		})
	}
}

type Minimum float64

func (m Minimum) Validate(path string, value interface{}, errs *[]Error) {
	val,ok:=value.(float64)
	if !ok{
		return
	}
	if val<float64(m){
		*errs = append(*errs,Error{
			Path: path,
			Info: appendString("value must be >=",strconv.FormatFloat(float64(m),'f',-1,64)),
		})
	}
}


type Enums []interface{}

func (enums Enums) Validate(path string, value interface{}, errs *[]Error) {
	if value == nil{
		return
	}
	for _, e := range enums {
		if e == value{
			return
		}
	}
	*errs = append(*errs,Error{
		Path: path,
		Info: fmt.Sprintf("must be one of %v",enums),
	})
}

type Required []string

func (r Required) Validate(path string, value interface{}, errs *[]Error) {
	m,ok:=value.(map[string]interface{})
	if !ok{
		return
	}
	for _,key := range r {
		if _,ok:=m[key];!ok{
			*errs = append(*errs,Error{
				Path: appendString(path,".",key),
				Info: "field is required",
			})
		}
	}
}



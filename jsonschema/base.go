package jsonschema

import (
	"fmt"
	"reflect"
	"strconv"
)

type Type struct {
	Path string
	Val string
}

func (t Type) Validate(c *ValidateCtx, value interface{}) {
	switch t.Val {
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
			c.AddError(Error{
				Path: t.Path,
				Info: "type must be object",
			})
		}
	case "string":
		if _,ok:=value.(string);!ok{
			c.AddError(Error{
				Path: t.Path,
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
			c.AddError(Error{
				Path: t.Path,
				Info: "type must be number",
			})
		}
	case "boolean","bool":
		if _,ok:=value.(bool);!ok{
			c.AddError(Error{
				Path: t.Path,
				Info: "type must be boolean",
			})
		}
	case "array":
		if _,ok:=value.([]interface{});!ok{
			c.AddError(Error{
				Path: t.Path,
				Info: "type must be array",
			})
		}
	}
}

func NewType(i interface{},path string,parent Validator)(Validator,error){
	iv, ok := i.(string)
	if !ok {
		return nil, fmt.Errorf("value of 'type' must be string! v:%v", i)
	}
	return Type{
		Val: iv,
		Path: path,
	},nil
}



type MaxLength struct {
	Val int
	Path string
}

func (l *MaxLength)Validate(c *ValidateCtx,value interface{}){

	switch value.(type) {
	case string:
		if len(value.(string)) > int(l.Val){
			c.AddError(Error{
				Info: "length must be <= "+strconv.Itoa(int(l.Val)),
			})
		}
	case []interface{}:
		if len(value.([]interface{})) >int(l.Val){
			c.AddError(Error{
				Info: "length must be <= "+strconv.Itoa(int(l.Val)),
			})
		}
	}

}

func NewMaxLen(i interface{},path string,parent Validator) (Validator, error) {
	v, ok := i.(float64)
	if !ok {
		return nil, fmt.Errorf("value of 'maxLength' must be int: %v", i)
	}
	return &MaxLength{
		Path: path,
		Val: int(v),
	}, nil
}

func NewMinLen(i interface{},path string,parent Validator) (Validator, error) {
	v, ok := i.(float64)
	if !ok {
		return nil, fmt.Errorf("value of 'minLengtg' must be int: %v", i)
	}
	return &MinLength{
		Val: int(v),
		Path: path,
	}, nil
}

func NewMaximum(i interface{},path string,parent Validator) (Validator, error) {
	v, ok := i.(float64)
	if !ok {
		return nil, fmt.Errorf("value of 'maximum' must be int")
	}
	return &Maximum{
		Val: v,
		Path:path,
	}, nil
}

func NewMinimum(i interface{},path string,parent Validator) (Validator, error) {
	v, ok := i.(float64)
	if !ok {
		return nil, fmt.Errorf("value of 'minimum' must be int")
	}
	return  &Minimum{
		Path: path,
		Val: v,
	}, nil
}

type MinLength struct {
	Val int
	Path string
}

func (l *MinLength)Validate(c *ValidateCtx,value interface{}){
	switch value.(type) {
	case string:
		if len(value.(string)) < int(l.Val){
			c.AddError(Error{
				Info: "length must be >= "+strconv.Itoa(int(l.Val)),
			})
		}
	case []interface{}:
		if len(value.([]interface{})) <int(l.Val){
			c.AddError(Error{
				Info: "length must be >= "+strconv.Itoa(int(l.Val)),
			})
		}
	}
}

type Maximum struct {
	Val float64
	Path string
}

func (m *Maximum) Validate(c *ValidateCtx,value interface{}) {
	val,ok:=value.(float64)
	if !ok{
		return
	}
	if val >float64(m.Val){
		c.AddError(Error{
			Info: appendString("value must be <=",strconv.FormatFloat(float64(m.Val),'f',-1,64)),
		})
	}
}

type Minimum struct {
	Val float64
	Path string
}

func (m Minimum) Validate(c *ValidateCtx,value interface{}) {
	val,ok:=value.(float64)
	if !ok{
		return
	}
	if val<(m.Val){
		c.AddError(Error{
			Info: appendString("value must be >=",strconv.FormatFloat(m.Val,'f',-1,64)),
		})
	}
}


type Enums struct {
	Val []interface{}
	Path string
}

func (enums *Enums) Validate(c *ValidateCtx,value interface{}) {
	if value == nil{
		return
	}
	for _, e := range enums.Val {
		if e == value{
			return
		}
	}

	for _,e:= range enums.Val{
		if Equal(e,value){
			return
		}
	}
	c.AddError(Error{
		Info: fmt.Sprintf("must be one of %v",enums),
	})
}

func NewEnums(i interface{},path string,parent Validator) (Validator, error) {
	arr, ok := i.([]interface{})
	if !ok {
		return nil, fmt.Errorf("value of 'enums' must be arr:%v", arr)
	}
	return &Enums{
		Val: arr,
		Path: path,
	}, nil
}


type Required struct {
	Val []string
	Path string
}

func (r Required) Validate(c *ValidateCtx,value interface{}) {
	m,ok:=value.(map[string]interface{})
	if !ok{
		return
	}
	for _,key := range r.Val {
		if _,ok:=m[key];!ok{
			c.AddError(Error{
				Path: appendString(r.Path,".",key),
				Info: "field is required",
			})
		}
	}
}


func NewRequired(i interface{},path string,parent Validator) (Validator, error) {
	arr, ok := i.([]interface{})
	if !ok {
		return nil, fmt.Errorf("value of 'required' must be array:%v", i)
	}
	req := make([]string, len(arr))
	for idx, item := range arr {
		itemStr, ok := item.(string)
		if !ok {
			return nil, fmt.Errorf("value of 'required item' must be string:%v of %v", item, i)
		}
		req[idx] = itemStr
	}
	return &Required{
		Val: req,
		Path: path,
	}, nil
}

type Items struct {
	Val ArrProp
	Path string
}

func (i Items) Validate(c *ValidateCtx,value interface{}) {
	if value == nil {
		return
	}
	arr, ok := value.([]interface{})
	if !ok {
		return
	}
	for _, item := range arr {
		for _, validator := range i.Val.Val {
			if validator.Val != nil {
				validator.Val.Validate(c,item)
			}
		}
	}
}
func NewItems(i interface{},path string, parent Validator) (Validator, error) {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("cannot create items with not object type: %v", i)
	}
	p, err := NewProp(m,path)
	if err != nil {
		return nil, err
	}
	return &Items{
		Val: p.(ArrProp),
		Path: path+"[?]",
	}, nil
}

package jsonschema

import (
	"fmt"
	"strconv"
)
// 性能提升

func init(){
	ShowCompletePath = false
}

type Properties2 struct {
	properties map[string]Validator
	constVals map[string]*ConstVal
	defaultVals map[string]*DefaultVal
	replaceKeys map[string]ReplaceKey
}

var ShowCompletePath bool

func (p *Properties2) Validate(path string, value interface{}, errs *[]Error) {
	if value == nil{
		return
	}
	if m,ok:=value.(map[string]interface{});ok {
		for k, v := range m {
			pv:=p.properties[k]
			if pv == nil{
				*errs = append(*errs,Error{
					Path: appendString(path,".",k),
					Info: "unknown field",
				})
				continue
			}
			if ShowCompletePath{
				pv.Validate(appendString(path,".",k),v,errs)

			}else{
				pv.Validate(k,v,errs)
			}
		}

		for key, val := range p.constVals {
			m[key] = val.Val
		}

		for key, val := range p.defaultVals {
			if _,ok:=m[key];!ok{
				m[key] = val.Val
			}
		}
		for key, rpk := range p.replaceKeys {
			if mv,ok:= m[key];ok{
				_,exist:=m[string(rpk)]
				if exist{ // 如果要替换的key 已经存在，不替换
					continue
				}
				m[string(rpk)] = mv

			}
		}

	}
	//else{
	//	*errs = append(*errs,Error{
	//		Path: path,
	//		Info: "type is not object",
	//	})
	//}
}

type Properties map[string]Validator

func (p Properties) Validate(path string, value interface{}, errs *[]Error) {
	if value == nil{
		return
	}
	if m,ok:=value.(map[string]interface{});ok{
		for k, v := range m {
			pv:=p[k]
			if pv == nil{
				*errs = append(*errs,Error{
					Path: appendString(path,".",k),
					Info: "unknown field",
				})
				continue
			}
			// p.set(path)
			//p.next()
			pv.Validate(appendString(path,".",k),v,errs)
			//p.back()
		}
		//return

		//设置默认值和定植
		for pk, vad := range p {
			prop,ok:=vad.(ArrProp)
			if !ok{
				continue
			}
			constVal:=prop.Get("constVal")
			if constVal != nil{
				cvl:=constVal.(*ConstVal)
				m[pk] = cvl.Val
			}

			defaultVal:=prop.Get("defaultVal")
			if defaultVal != nil && m[pk]==nil{
				dfv:=defaultVal.(*DefaultVal)
				m[pk] = dfv.Val
			}
		}

		// 设置替换的key
		for k, v := range m {
			if prop,ok:=p[k].(ArrProp);ok{
				replaceKey := prop.Get("replaceKey")
				if replaceKey != nil{
					rpk:=replaceKey.(ReplaceKey)
					m[string(rpk)] = v
				}
			}
		}

	}else{
		*errs = append(*errs,Error{
			Path: path,
			Info: "type is not object",
		})
	}

}




type Prop map[string]Validator

func (p Prop) Validate(path string, value interface{}, errs *[]Error) {
	for _, v := range p {
		v.Validate(path,value,errs)
	}
}

func (p Prop)Get(key string)Validator{
	return p[key]
}


type Items ArrProp

func (i Items) Validate(path string, value interface{}, errs *[]Error) {
	if value == nil{
		return
	}
	arr,ok:=value.([]interface{})
	if !ok{
		return
	}
	for idx, item := range arr {
		for _, validator := range i {
			if validator.Val != nil{
				validator.Val.Validate(appendString(path,"[",strconv.Itoa(idx),"]"),item,errs)
			}
		}
	}
}
func NewItems(i interface{},parent Validator) (Validator, error) {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("cannot create items with not object type: %v", i)
	}
	p ,err:= NewProp(m)
	if err != nil{
		return nil,err
	}
	return Items(p.(ArrProp)), nil
}
type PropItem struct {
	Key string
	Val Validator
}

// 数组的for 比  map for快很多,在数据不大的情况下，for 也比 map get  快
type ArrProp []PropItem

func (a ArrProp) Validate(path string, value interface{}, errs *[]Error) {
	for _, item := range a {
		if item.Val == nil{
			continue
		}
		item.Val.Validate(path,value,errs)
	}
}
func (a ArrProp)Get(key string)Validator{
	for _, item := range a {
		if item.Key == key{
			return item.Val
		}
	}
	return nil
}

type FlexProperties map[string]Validator

func (f FlexProperties) Validate(path string, value interface{}, errs *[]Error) {
	m,ok:=value.(map[string]interface{})
	if !ok{
		return
	}
	for key, val := range m {
		vad:=f[key]
		if vad == nil{
			continue
		}
		f[key].Validate(appendString(path,".",key),val,errs)
	}
}
// 宽松校验器，允许存在不在properties 中的值
func NewFlexProperties(i interface{},parent Validator)(Validator,error){
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("cannot create flexProperties with not object type: %v", i)
	}
	p := FlexProperties{}
	for key, val := range m {
		vad, err := NewProp(val)
		if err != nil {
			return nil, err
		}
		p[key] = vad
	}
	return p,nil
}


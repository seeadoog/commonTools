package jsonschema

import (
	"fmt"
	"reflect"
	"strconv"
)

// 性能提升

func init() {
	ShowCompletePath = false
}

type Properties2 struct {
	properties  map[string]Validator
	constVals   map[string]*ConstVal
	defaultVals map[string]*DefaultVal
	replaceKeys map[string]ReplaceKey
}

var ShowCompletePath bool

func (p *Properties2) Validate(path string, value interface{}, errs *[]Error) {
	if value == nil {
		return
	}
	if m, ok := value.(map[string]interface{}); ok {
		for k, v := range m {
			pv := p.properties[k]
			if pv == nil {
				*errs = append(*errs, Error{
					Path: appendString(path, ".", k),
					Info: "unknown field",
				})
				continue
			}

			pv.Validate(appendString(path, ".", k), v, errs)

		}

		for key, val := range p.constVals {
			m[key] = val.Val
		}

		for key, val := range p.defaultVals {
			if _, ok := m[key]; !ok {
				m[key] = val.Val
			}
		}
		for key, rpk := range p.replaceKeys {
			if mv, ok := m[key]; ok {
				_, exist := m[string(rpk)]
				if exist { // 如果要替换的key 已经存在，不替换
					continue
				}
				m[string(rpk)] = mv

			}
		}
	} else {
		rv := reflect.ValueOf(value)
		switch rv.Kind() {
		case reflect.Ptr:
			if rv.IsNil() {
				return
			}
			rv = rv.Elem()
			fallthrough
		case reflect.Struct:
			rt := rv.Type()
			for i := 0; i < rv.NumField(); i++ {
				ft := rt.Field(i)
				propName := ft.Tag.Get("json")
				if propName == "" {
					propName = ft.Name
				}
				vad := p.properties[propName]
				if vad == nil {
					return
				}
				fv := rv.Field(i)
				if fv.CanInterface() {
					vad.Validate(propName, fv.Interface(), errs)
				}
				// set constVal ,struct 类型无法知道目标值是否为空，无法设定默认值
				var vv interface{} = nil
				constv := p.constVals[propName]
				if constv != nil {
					vv = constv.Val
				}
				if vv == nil{
					continue
				}
				setV := reflect.ValueOf(vv)
				if setV.Kind() == fv.Kind() {
					fv.Set(setV)
				} else if setV.Kind() == reflect.Float64 {
					switch fv.Kind() {
					case reflect.Int, reflect.Int32, reflect.Int64, reflect.Int8, reflect.Int16:
						fv.SetInt(int64(setV.Interface().(float64)))
					case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
						fv.SetUint(uint64(setV.Interface().(float64)))
					case reflect.Float32:
						fv.SetFloat(setV.Interface().(float64))
					}
				}

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

type Items ArrProp

func (i Items) Validate(path string, value interface{}, errs *[]Error) {
	if value == nil {
		return
	}
	arr, ok := value.([]interface{})
	if !ok {
		return
	}
	for idx, item := range arr {
		for _, validator := range i {
			if validator.Val != nil {
				validator.Val.Validate(appendString(path, "[", strconv.Itoa(idx), "]"), item, errs)
			}
		}
	}
}
func NewItems(i interface{}, parent Validator) (Validator, error) {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("cannot create items with not object type: %v", i)
	}
	p, err := NewProp(m)
	if err != nil {
		return nil, err
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
		if item.Val == nil {
			continue
		}
		item.Val.Validate(path, value, errs)
	}
}
func (a ArrProp) Get(key string) Validator {
	for _, item := range a {
		if item.Key == key {
			return item.Val
		}
	}
	return nil
}

type FlexProperties map[string]Validator

func (f FlexProperties) Validate(path string, value interface{}, errs *[]Error) {
	m, ok := value.(map[string]interface{})
	if !ok {
		return
	}
	for key, val := range m {
		vad := f[key]
		if vad == nil {
			continue
		}
		f[key].Validate(appendString(path, ".", key), val, errs)
	}
}

// 宽松校验器，允许存在不在properties 中的值
func NewFlexProperties(i interface{}, parent Validator) (Validator, error) {
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
	return p, nil
}

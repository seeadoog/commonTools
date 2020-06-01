package jsonschema

import "fmt"

func init() {
	// 这些显示放在funcs 里面时，不让编译通过，透。。。
	RegisterValidator("properties", NewProperties2)
	RegisterValidator("items", NewItems)
	RegisterValidator("anyOf", NewAnyOf)
	RegisterValidator("if", NewIf)
	RegisterValidator("else", NewElse)
	RegisterValidator("then", NewThen)
	RegisterValidator("flexProperties", NewFlexProperties)
	RegisterValidator("not", NewNot)
	RegisterValidator("allOf", NewAllOf)
	RegisterValidator("dependencies", NewDependencies)
	RegisterValidator("keyMatch", NewKeyMatch)
	RegisterValidator("setVal", NewSetVal)
	RegisterValidator("script",NewScript)
	RegisterValidator("switch",NewSwitch)

}
// 忽略的校验器
var ignoreKeys = map[string]int{
	"title":1,
	"comment":1,
}

var funcs = map[string]ValidateFunc{
	"type":       NewType,
	"maxLength":  NewMaxLen,
	"minLength":  NewMinLen,
	"maximum":    NewMaximum,
	"minimum":    NewMinimum,
	"required":   NewRequired,
	"constVal":   NewConstVal,
	"defaultVal": NewDefaultVal,
	"replaceKey": NewReplaceKey,
	"enums":      NewEnums,
	"enum":      NewEnums,
	"pattern":    NewPattern,
}

func RegisterValidator(name string, fun ValidateFunc) {
	if funcs[name] != nil {
		panicf("register validator error!%s already exists", name)
	}
	funcs[name] = fun
}

type ValidateFunc func(i interface{},parent Validator) (Validator, error)

func NewType(i interface{},parent Validator) (Validator, error) {
	iv, ok := i.(string)
	if !ok {
		return nil, fmt.Errorf("value of 'type' must be string! v:%v", i)
	}
	return Type(iv), nil
}

func NewMaxLen(i interface{},parent Validator) (Validator, error) {
	v, ok := i.(float64)
	if !ok {
		return nil, fmt.Errorf("value of 'maxLength' must be int: %v", i)
	}
	return MaxLength(v), nil
}

func NewMinLen(i interface{},parent Validator) (Validator, error) {
	v, ok := i.(float64)
	if !ok {
		return nil, fmt.Errorf("value of 'minLengtg' must be int: %v", i)
	}
	return MinLength(v), nil
}

func NewMaximum(i interface{},parent Validator) (Validator, error) {
	v, ok := i.(float64)
	if !ok {
		return nil, fmt.Errorf("value of 'maximum' must be int")
	}
	return Maximum(v), nil
}

func NewMinimum(i interface{},parent Validator) (Validator, error) {
	v, ok := i.(float64)
	if !ok {
		return nil, fmt.Errorf("value of 'minimum' must be int")
	}
	return Minimum(v), nil
}

func NewProp(i interface{}) (Validator, error) {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("cannot create prop with not object type: %v", i)
	}
	p := make(ArrProp, len(m))
	idx:=0
	for key, val := range m {
		if ignoreKeys[key] >0{
			continue
		}
		if funcs[key] == nil {
			return nil, fmt.Errorf("%s is unknown validator", key)
		}
		vad, err := funcs[key](val,p)
		if err != nil {
			return nil, err
		}
		//p[key] = vad
		p[idx] =  PropItem{Key: key, Val: vad}
		idx++

	}
	return p, nil
}

func NewProperties(i interface{},parent Validator) (Validator, error) {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("cannot create properties with not object type: %v", i)
	}
	p := Properties{}
	for key, val := range m {
		vad, err := NewProp(val)
		if err != nil {
			return nil, err
		}
		p[key] = vad
	}
	return p, nil

}

func NewProperties2(i interface{},parent Validator) (Validator, error) {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("cannot create properties with not object type: %v", i)
	}
	p := &Properties2{
		properties: map[string]Validator{},
		replaceKeys: map[string]ReplaceKey{},
		constVals: map[string]*ConstVal{},
		defaultVals: map[string]*DefaultVal{},
	}
	for key, val := range m {
		vad, err := NewProp(val)
		if err != nil {
			return nil, err
		}
		p.properties[key] = vad
	}

	for key, val := range p.properties {
		prop,ok:=val.(ArrProp)
		if !ok{
			continue
		}
		constVal,ok:=prop.Get("constVal").(*ConstVal)
		if ok{
			p.constVals[key] =constVal
		}
		defaultVal,ok:=prop.Get("defaultVal").(*DefaultVal)
		if ok{
			p.defaultVals[key] = defaultVal
		}
		replaceKey,ok:=prop.Get("replaceKey").(ReplaceKey)
		if ok{
			p.replaceKeys[key] = replaceKey
		}
	}

	return p, nil

}
func NewItems(i interface{},parent Validator) (Validator, error) {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("cannot create items with not object type: %v", i)
	}
	p := Items{}
	for key, val := range m {
		vad, err := NewProp(val)
		if err != nil {
			return nil, err
		}
		p[key] = vad
	}
	return p, nil
}

func NewRequired(i interface{},parent Validator) (Validator, error) {
	arr, ok := i.([]interface{})
	if !ok {
		return nil, fmt.Errorf("value of 'required' must be array:%v", i)
	}
	req := make(Required, len(arr))
	for idx, item := range arr {
		itemStr, ok := item.(string)
		if !ok {
			return nil, fmt.Errorf("value of 'required item' must be string:%v of %v", item, i)
		}
		req[idx] = itemStr
	}
	return req, nil
}

func NewEnums(i interface{},parent Validator) (Validator, error) {
	arr, ok := i.([]interface{})
	if !ok {
		return nil, fmt.Errorf("value of 'enums' must be arr:%v", arr)
	}
	return Enums(arr), nil
}



package jsonschema

import (
	"fmt"
)

type AnyOf []Validator

func (a AnyOf) Validate(path *pathTree, value interface{}, errs *[]Error) {
	allErrs:=[]Error{}
	for _, validator := range a {
		e:=[]Error{}
		validator.Validate(path,value,&e)
		if len(e) == 0{
			return
		}
		allErrs = append(allErrs,e...)
	}
	// todo 区分errors
	*errs = append(*errs,allErrs...)
}

func NewAnyOf( i interface{},parent Validator)(Validator,error){
	m,ok:=i.([]interface{})
	if !ok{
		return nil,fmt.Errorf("value of anyOf must be array:%v",i)
	}
	any:=AnyOf{}
	for idx, v := range m {
		ip,err:=NewProp(v)
		if err != nil{
			return nil, fmt.Errorf("anyOf index:%d is invalid:%w %v",idx,err,v)
		}
		any = append(any,ip)
	}
	return any,nil
}

type If struct {
	parent ArrProp
	v Validator
}

func (i If) Validate(path *pathTree, value interface{}, errs *[]Error) {
	ifErrs:=[]Error{}
	i.v.Validate(path,value,&ifErrs)
	if len(ifErrs)==0{
		if i.parent!= nil{
			then,ok:=i.parent.Get("then").(Then)
			if ok{
				then.v.Validate(path,value,errs)
			}
		}
	}else{
		if i.parent!= nil{
			elsev,ok:=i.parent.Get("else").(Else)
			if ok{
				elsev.v.Validate(path,value,errs)
			}
		}
	}
}

func NewIf(i interface{},parent Validator)(Validator,error){
	ifp,err:=NewProp(i)
	if err != nil{
		return nil, err
	}
	iff:=&If{
		v: ifp,
	}
	pp,ok:=parent.(ArrProp)
	if ok{
		iff.parent = pp
	}
	return iff,nil
}
type Then struct {
	v Validator
}

func (t Then) Validate(path *pathTree, value interface{}, errs *[]Error) {
	// then 不能主动调用
}

type Else struct {
	v Validator
}

func (e Else) Validate(path *pathTree, value interface{}, errs *[]Error) {
	//panic("implement me")
}

func NewThen(i interface{},parent Validator)(Validator,error){
	v,err:= NewProp(i)
	if err != nil{
		return nil, err
	}
	return Then{
		v: v,
	},nil
}

func NewElse(i interface{},parent Validator)(Validator,error){
	v,err:= NewProp(i)
	if err != nil{
		return nil, err
	}
	return Else{
		v: v,
	},nil
}

type Not struct {
	v Validator
}

func (n Not) Validate(path *pathTree, value interface{}, errs *[]Error) {
	ners:=[]Error{}
	n.v.Validate(path,value,&ners)
	//fmt.Println(ners,value)
	if len(ners) ==0{
		*errs = append(*errs,Error{
			Path: path.String(),
			Info: "is not valid",
		})
	}
}

func NewNot(i interface{},parent Validator)(Validator,error){
	p,err:=NewProp(i)
	if err != nil{
		return nil, err
	}
	return Not{v: p},nil
}

type AllOf []Validator

func (a AllOf) Validate(path *pathTree, value interface{}, errs *[]Error) {
	for _, validator := range a {
		validator.Validate(path,value,errs)
	}
}

func NewAllOf(i interface{},parent Validator)(Validator,error){
	arr,ok:=i.([]interface{})
	if !ok{
		return nil,fmt.Errorf("value of 'allOf' must be array: %v",i)
	}
	all:=AllOf{}
	for _, ai := range arr {
		iv,err:=NewProp(ai)
		if err != nil{
			return nil, err
		}
		all = append(all,iv)
	}
	return all,nil
}

type Dependencies map[string][]string

func (d Dependencies) Validate(path *pathTree, value interface{}, errs *[]Error) {
	m,ok:=value.(map[string]interface{})
	if !ok{
		return
	}
	// 如果存在key，那么必须存在某些key
	for key, vals := range d {
		_,ok:=m[key]
		if ok{
			for _, val := range vals {
				_,ok = m[val]
				if !ok{
					*errs = append(*errs,Error{
						Path: appendString(path.String(),".",val),
						Info: "is required",
					})
				}
			}
		}
	}
}

func NewDependencies(i interface{},parent Validator)(Validator,error){
	m,ok:=i.(map[string]interface{})
	if !ok{
		return nil, fmt.Errorf("value of dependencies must be map[string][]string :%v", i)
	}
	vad:=Dependencies{}
	for key, arris := range m {
		arrs,ok:=arris.([]interface{})
		if !ok{
			return nil, fmt.Errorf("value of dependencies must be map[string][]string :%v", i)
		}
		strs:=make([]string, len(arrs))
		for idx, item := range arrs {
			str,ok:=item.(string)
			if !ok{
				return nil, fmt.Errorf("value of dependencies must be map[string][]string :%v", i)

			}
			strs[idx] = str
		}
		vad[key] = strs

	}
	return vad, nil
}

/*
{
	"keyMatch":{
		"key1":"biaoge"
	}
}
 */

type KeyMatch map[string]interface{}

func (k KeyMatch) Validate(path *pathTree, value interface{}, errs *[]Error) {
	m,ok:=value.(map[string]interface{})
	if !ok{
		*errs = append(*errs,Error{
			Path: path.String(),
			Info: "type is not object",
		})
	}
	for key, want := range k {
		target:=m[key]
		if target != want {
			*errs = append(*errs,Error{
				Path: appendString(path.String(),".",key),
				Info: fmt.Sprintf("value must be %v",want),
			})
		}
	}
}

func NewKeyMatch(i interface{},parent Validator)(Validator,error){
	m,ok:=i.(map[string]interface{})
	if !ok{
		return nil,fmt.Errorf("value of keyMatch must be map[string]interface{} :%v",i)
	}
	return KeyMatch(m),nil
}

/*
{
	"switch":{
		"key":"key1",
		"case":{
			"v1":{},
			"v2":{}
		},
		"default":{

		}
	}
}
 */
type Switch struct {
	Switch string
	Case map[string]Validator
	Default Validator
}

func (s Switch) Validate(path *pathTree, value interface{}, errs *[]Error) {
	m,ok:=value.(map[string]interface{})
	if !ok{
		if s.Default != nil{
			s.Default.Validate(path,value,errs)
		}
		return
	}
	for cas, validator := range s.Case {
		if cas ==String(m[s.Switch]){
			validator.Validate(path,value,errs)
			return
		}
	}
	if s.Default != nil{
		s.Default.Validate(path,value,errs)
	}
}

func NewSwitch(i interface{},parent Validator)(Validator,error){
	m,ok:=i.(map[string]interface{})
	if !ok{
		return nil, fmt.Errorf("value of Switch must be map :%v", i)
	}

	swth,ok:=m["key"].(string)
	if !ok{
		return nil, fmt.Errorf("switch key must be string:%v",i)
	}
	s:=&Switch{
		Switch: swth,
		Case: map[string]Validator{},
	}
	cases,ok:=m["cases"].(map[string]interface{})
	if !ok{
		return nil, fmt.Errorf("cases must be map:%v",i)
	}
	for key, val := range cases {
		vad,err:=NewProp(val)
		if err != nil{
			return nil, err
		}
		s.Case[key] = vad
	}
	def:=m["default"]
	if def != nil{
		defv,err:=NewProp(def)
		if err != nil{
			return nil, err
		}
		s.Default = defv
	}
	return s, nil
}
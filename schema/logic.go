package schema

import (
	"fmt"
)

type AnyOf []Validator

func (a AnyOf) Validate(c *ValidateCtx,value interface{}) {
	allErrs:=[]Error{}
	for _, validator := range a {
		cb:=c.Clone()
		validator.Validate(cb,value)
		if len(cb.errors) == 0{
			return
		}
		allErrs = append(allErrs,cb.errors...)
	}
	// todo 区分errors

	c.AddErrors(allErrs...)
}

func NewAnyOf( i interface{},path string,parent Validator)(Validator,error){
	m,ok:=i.([]interface{})
	if !ok{
		return nil,fmt.Errorf("value of anyOf must be array:%v",i)
	}
	any:=AnyOf{}
	for idx, v := range m {
		ip,err:=NewProp(v,path)
		if err != nil{
			return nil, fmt.Errorf("anyOf index:%d is invalid:%w %v",idx,err,v)
		}
		any = append(any,ip)
	}
	return any,nil
}

type If struct {
	parent *ArrProp
	v Validator
}

func (i If) Validate(c *ValidateCtx,value interface{}) {
	cif:=c.Clone()
	i.v.Validate(cif,value)
	if len(cif.errors)==0{
		if i.parent!= nil{
			then,ok:=i.parent.Get("then").(Then)
			if ok{
				then.v.Validate(c,value)
			}
		}
	}else{
		if i.parent!= nil{
			elsev,ok:=i.parent.Get("else").(Else)
			if ok{
				elsev.v.Validate(c,value)
			}
		}
	}
}

func NewIf(i interface{},path string,parent Validator)(Validator,error){
	ifp,err:=NewProp(i,path)
	if err != nil{
		return nil, err
	}
	iff:=&If{
		v: ifp,
	}
	pp,ok:=parent.(*ArrProp)
	if ok{
		iff.parent = pp
	}
	return iff,nil
}
type Then struct {
	v Validator
}

func (t Then) Validate(c *ValidateCtx,value interface{}) {
	// then 不能主动调用
}

type Else struct {
	v Validator
}

func (e Else) Validate(c *ValidateCtx,value interface{}) {
	//panic("implement me")
}

func NewThen(i interface{},path string,parent Validator)(Validator,error){
	v,err:= NewProp(i,path)
	if err != nil{
		return nil, err
	}
	return Then{
		v: v,
	},nil
}

func NewElse(i interface{},path string,parent Validator)(Validator,error){
	v,err:= NewProp(i,path)
	if err != nil{
		return nil, err
	}
	return Else{
		v: v,
	},nil
}

type Not struct {
	v Validator
	Path string
}

func (n Not) Validate(c *ValidateCtx,value interface{}) {
	cn:=c.Clone()
	n.v.Validate(cn,value)
	//fmt.Println(ners,value)
	if len(cn.errors) ==0{
		c.AddErrors(Error{
			Path: n.Path,
			Info: "is not valid",
		})
	}
}

func NewNot(i interface{},path string,parent Validator)(Validator,error){
	p,err:=NewProp(i,path)
	if err != nil{
		return nil, err
	}
	return Not{v: p},nil
}

type AllOf []Validator

func (a AllOf) Validate(c *ValidateCtx,value interface{}) {
	for _, validator := range a {
		validator.Validate(c,value)
	}
}

func NewAllOf(i interface{},path string,parent Validator)(Validator,error){
	arr,ok:=i.([]interface{})
	if !ok{
		return nil,fmt.Errorf("value of 'allOf' must be array: %v",i)
	}
	all:=AllOf{}
	for _, ai := range arr {
		iv,err:=NewProp(ai,path)
		if err != nil{
			return nil, err
		}
		all = append(all,iv)
	}
	return all,nil
}

type Dependencies struct {
	Val map[string][]string
	Path string
}

func (d *Dependencies) Validate(c *ValidateCtx,value interface{}) {
	m,ok:=value.(map[string]interface{})
	if !ok{
		return
	}
	// 如果存在key，那么必须存在某些key
	for key, vals := range d.Val{
		_,ok:=m[key]
		if ok{
			for _, val := range vals {
				_,ok = m[val]
				if !ok{
					c.AddErrors(Error{
						Path: appendString(d.Path,".",val),
						Info: "is required",
					})
				}
			}
		}
	}
}

func NewDependencies(i interface{},path string,parent Validator)(Validator,error){
	m,ok:=i.(map[string]interface{})
	if !ok{
		return nil, fmt.Errorf("value of dependencies must be map[string][]string :%v", i)
	}
	vad:=&Dependencies{
		Val: map[string][]string{},
		Path: path,
	}
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
		vad.Val[key] = strs

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

type KeyMatch struct {
	Val map[string]interface{}
	Path string
}

func (k *KeyMatch) Validate(c *ValidateCtx,value interface{}) {
	m,ok:=value.(map[string]interface{})
	if !ok{
		c.AddError(Error{
			Path: k.Path,
			Info: "val is not object",
		})
	}
	for key, want := range k.Val {
		target:=m[key]
		if target != want {
			c.AddError(Error{
				Path: appendString(k.Path,".",key),
				Info: fmt.Sprintf("value must be %v",want),
			})
		}
	}
}

func NewKeyMatch(i interface{},path string,parent Validator)(Validator,error){
	m,ok:=i.(map[string]interface{})
	if !ok{
		return nil,fmt.Errorf("value of keyMatch must be map[string]interface{} :%v",i)
	}
	return &KeyMatch{
		Val: m,
		Path: path,
	},nil
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

func (s Switch) Validate(c *ValidateCtx,value interface{}) {
	m,ok:=value.(map[string]interface{})
	if !ok{
		if s.Default != nil{
			s.Default.Validate(c,value)
		}
		return
	}
	for cas, validator := range s.Case {
		if cas ==String(m[s.Switch]){
			validator.Validate(c,value)
			return
		}
	}
	if s.Default != nil{
		s.Default.Validate(c,value)
	}
}

func NewSwitch(i interface{},path string,parent Validator)(Validator,error){
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
		vad,err:=NewProp(val,path)
		if err != nil{
			return nil, err
		}
		s.Case[key] = vad
	}
	def:=m["default"]
	if def != nil{
		defv,err:=NewProp(def,path)
		if err != nil{
			return nil, err
		}
		s.Default = defv
	}
	return s, nil
}
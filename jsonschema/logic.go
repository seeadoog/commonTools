package jsonschema

import (
	"fmt"
)

type AnyOf []Validator

func (a AnyOf) Validate(path string, value interface{}, errs *[]Error) {
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

func (i If) Validate(path string, value interface{}, errs *[]Error) {
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
				elsev.Validate(path,value,errs)
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

func (t Then) Validate(path string, value interface{}, errs *[]Error) {
	// then 不能主动调用
}

type Else struct {
	v Validator
}

func (e Else) Validate(path string, value interface{}, errs *[]Error) {
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

func (n Not) Validate(path string, value interface{}, errs *[]Error) {
	ners:=[]Error{}
	n.v.Validate(path,value,&ners)
	//fmt.Println(ners,value)
	if len(ners) ==0{
		*errs = append(*errs,Error{
			Path: path,
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

func (a AllOf) Validate(path string, value interface{}, errs *[]Error) {
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


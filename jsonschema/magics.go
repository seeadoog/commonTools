package jsonschema

import "fmt"

type ConstVal struct {
	Val interface{}
}

func (c ConstVal) Validate(path string, value interface{}, errs *[]Error) {

}

type DefaultVal struct {
	Val interface{}
}

func (d DefaultVal) Validate(path string, value interface{}, errs *[]Error) {

}

type ReplaceKey string

func (r ReplaceKey) Validate(path string, value interface{}, errs *[]Error) {

}

func NewConstVal(i interface{},parent Validator) (Validator, error) {
	return &ConstVal{
		Val: i,
	}, nil
}

func NewDefaultVal(i interface{},parent Validator) (Validator, error) {
	return &DefaultVal{i}, nil
}

func NewReplaceKey(i interface{},parent Validator) (Validator, error) {
	s, ok := i.(string)
	if !ok {
		return nil, fmt.Errorf("value of 'replaceKey' must be string :%v", i)
	}
	return ReplaceKey(s), nil

}

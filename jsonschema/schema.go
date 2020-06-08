package jsonschema

import (
	"encoding/json"
	"errors"
	"strings"
)

type Schema struct {
	prop Validator
	i interface{}
}

func (s *Schema)UnmarshalJSON(b []byte)error{
	var i interface{}
	if err:=json.Unmarshal(b,&i);err != nil{
		return err
	}
	s.i = i
	p ,err := NewProp(i)
	if err != nil{
		return err
	}
	s.prop = p
	return nil
}

func (s *Schema)MarshalJSON()(b []byte,err error){
	data,err:=json.Marshal(s.i)
	if err != nil{
		return nil,err
	}
	return data,nil

}

func (s *Schema)Validate(i interface{})error{
	errs:=[]Error{}
	path:=newPathTree()
	path.path = "$"
	s.prop.Validate(path,i,&errs)
	if len(errs) == 0{
		return nil
	}
	pool.Put(path)
	return errors.New(errsToString(errs))
}



func errsToString(errs []Error)string{
	sb:=strings.Builder{}
	for _, err := range errs {
		sb.WriteString(appendString("'",err.Path,"' ",err.Info,"; "))
	}
	return sb.String()
}
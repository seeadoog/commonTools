package simple_json

import "encoding/json"

type Element struct {
	data interface{}
}

func (e *Element)Get(key string)*Element{
	if e== nil{
		return nil
	}
	m,ok:=e.data.(map[string]interface{})
	if ok{
		return &Element{
			data:m[key],
		}
	}
	return nil
}

func (e *Element)UnmarshalJSON(b []byte)error{
	return json.Unmarshal(b,&e.data)
}

func (e *Element)MarshalJSON()([]byte,error){
	if e== nil{
		return nil,nil
	}
	b,err:=json.Marshal(e.data)
	return b,err
}
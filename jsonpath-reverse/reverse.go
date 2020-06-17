package jsonpath_reverse

import (
	"errors"
	"strconv"
	"strings"
	"fmt"
	"sync"
)

type QueryProp struct {
	Query string
	Value interface{}
}

type SwitchExp struct {
	SrcExp  string
	DataExp string
}

const (
	TYPE_KEY = -1
)

// src must be map[string]interface{}
// Marshal() do not support expression start with  $,$[0] .
//if you need this function ,please use MarshalInterface()
func Marshal(query string, src interface{}, value interface{}) error {
	return marshal(query, src, value, 0)
}

func marshal(query string, src interface{}, value interface{}, start int) error {
	tks,err:=parseExps(query)
	if err != nil {
		return err
	}
	var cp = src
	return parserToken(tks, cp, value)
}

func Marshals(querys []QueryProp) (tmp map[string]interface{}, err error) {
	m := map[string]interface{}{}
	for _, v := range querys {
		if err := Marshal(v.Query, m, v.Value); err != nil {
			return nil, err
		}
	}
	return m, nil

}

//resolve the token is a field or array
//return -1 and field name if token is field
//return the idx of token and array name if token is array
func yyp(token string) (string, int, error) {
	numidx_start := 0
	numidx_end := 0
	for k, v := range token {
		t := string(v)
		if t == "[" {
			numidx_start = k
		}
		if t == "]" {
			numidx_end = k
		}
	}
	if numidx_end > 0 && numidx_start >= 0 {
		num, err := strconv.Atoi(token[numidx_start+1 : numidx_end])
		if err != nil {
			return "", TYPE_KEY, err
		}
		return token[:numidx_start], num, nil
	}
	return token, TYPE_KEY, nil
}



var expCache = sync.Map{}

type token struct {
	idx int
	key string
}

type RefExps struct {
	isroot bool
	tokens []*token
}

func (exps *RefExps)UnmarshalJSON(b []byte)error{
	exp:=strings.Trim(string(b),"\"")
	pex,err:= parseExps(exp)
	if err != nil{
		return fmt.Errorf("parse ref error :%s,%w",exp,err)
	}
	exps.isroot = pex.isroot
	exps.tokens = pex.tokens
	return nil
}




func (exps *RefExps)Set(dst,value interface{})error{
	var cp = dst
	if cpi, ok := dst.(*interface{}); ok {
		if exps.isroot {
			*cpi = value
			return nil
		}
		if _, ok := (*cpi).(map[string]interface{}); ok {
			// do not handle
		} else if _, ok := (*cpi).([]interface{}); ok {
			cp = cpi
			goto done
		} else {
			//yp, idx, err := yyp(tks[0])
			//if err != nil {
			//	return err
			//}
			tkn:=exps.tokens[0]
			if tkn.idx == TYPE_KEY {
				*cpi = map[string]interface{}{}
			} else {
				if tkn.key != "" {
					*cpi = map[string]interface{}{}
				} else {
					*cpi = make([]interface{}, tkn.idx+1)
					cp = cpi
					goto done
				}

			}
		}
		cp = *cpi

	}
done:
	return parserToken(exps, cp, value)
}

func parseExps(query string)(*RefExps,error){
	tkns,ok:= expCache.Load(query)
	if ok{
		return tkns.(*RefExps),nil
	}
	rawtkns:=strings.Split(strings.TrimLeft(query, "$."), ".")
	parsedTkns:=make([]*token,0, len(rawtkns))
	for _, rawtkn := range rawtkns {
		key,idx,err:=yyp(rawtkn)
		if err !=nil{
			return nil,err
		}
		parsedTkns = append(parsedTkns,&token{
			idx:    idx,
			key:    key,
		})
	}

	exp:=&RefExps{
		tokens:parsedTkns,
	}
	if strings.TrimLeft(query, "$.")==""{
		exp.isroot = true
	}
	expCache.Store(query,exp)
	return exp,nil

}

//cahce tokens
var tokensCache sync.Map

func tokenize2(query string) ([]string, error) {
	tkn,ok:=tokensCache.Load(query)
	if !ok{
		tkns :=strings.Split(strings.TrimLeft(query, "$."), ".")
		tokensCache.Store(query,tkns)
		return tkns,nil
	}
	return tkn.([]string),nil
}

type CompiledRefExps struct {

}

//marshal and set the value to interface{}
//MarshalInterface() is power than Marshal().
// MarshalInterface() support expression such as $ ,$[0] which  Marshal() doesn't support
//attention that $dst must be *interface{}
func MarshalInterface(query string, dst interface{}, value interface{}) error {
	return marshalInterface(query, dst, value)
}




func marshalInterface(query string, dst interface{}, value interface{}) error {
	//fmt.Println(reflect.TypeOf(dst))
	exps, err := parseExps(query)
	if err != nil {
		return err
	}
	var cp = dst
	if cpi, ok := dst.(*interface{}); ok {
		if exps.isroot {
			*cpi = value
			return nil
		}
		if _, ok := (*cpi).(map[string]interface{}); ok {
			// do not handle
		} else if _, ok := (*cpi).([]interface{}); ok {
			cp = cpi
			goto done
		} else {
			//yp, idx, err := yyp(tks[0])
			//if err != nil {
			//	return err
			//}
			tkn:=exps.tokens[0]
			if tkn.idx == TYPE_KEY {
				*cpi = map[string]interface{}{}
			} else {
				if tkn.key != "" {
					*cpi = map[string]interface{}{}
				} else {
					*cpi = make([]interface{}, tkn.idx+1)
					cp = cpi
					goto done
				}

			}
		}
		cp = *cpi

	}
done:
	return parserToken(exps, cp, value)

}

func parserToken(tks *RefExps, cp, value interface{}) error {
	for k, v := range tks.tokens {
		//field, idx, err := yyp(v)
		//if err != nil {
		//	return err
		//}
		//
		field:=v.key
		idx:=v.idx
		if idx == TYPE_KEY {
			cpm, ok := cp.(map[string]interface{})
			if !ok {
				return errors.New(fmt.Sprintf("create field failed ,%s->parent cannot convert_ to map", field))
			}

			if k < len(tks.tokens)-1 {
				if cpm[field] == nil {
					cpm[field] = map[string]interface{}{}
				}
				cpm, ok = cpm[field].(map[string]interface{})
				if !ok {
					return errors.New(fmt.Sprintf("create field failed ,%s cannot convert_ to map", field))
				}
				cp = cpm
			} else {
				//filed is last token ,set value to interface
				cpm[field] = value
			}
		} else { //array
			if field == "" && k == 0 {  //root array
				cpi, ok := cp.(*interface{})
				//	fmt.Println(reflect.TypeOf(cp))
				if !ok {
					return errors.New("root is not pointer")
				}
				if _, ok := (*cpi).([]interface{}); !ok {
					return errors.New("root is not array")
				}
				if len((*cpi).([]interface{})) < idx+1 {
					for i := len((*cpi).([]interface{})); i < idx+1; i++ {
						*cpi = append((*cpi).([]interface{}), nil)
					}
				}

				if k < len(tks.tokens)-1 {
					for i := 0; i < idx+1; i++ {
						if (*cpi).([]interface{})[i] == nil {
							(*cpi).([]interface{})[i] = map[string]interface{}{}
						}
					}
					cp = (*cpi).([]interface{})[idx]
				} else {
					//filed is last token ,set value to interface
					(*cpi).([]interface{})[idx] = value
				}
				//fmt.Println((*cpi).([]interface{}))
				continue
			}

			cpm, ok := cp.(map[string]interface{})
			if !ok {
				return errors.New("nil array child")
			}
			if cpm[field] == nil {
				cpm[field] = make([]interface{}, idx+1)
			}
			cps, ok := cpm[field].([]interface{})
			if !ok {
				return errors.New(fmt.Sprintf("create array failed ,%s cannot convert2 to array", field))
			}
			lenmap := len(cps)
			if lenmap < idx+1 {
				for i := lenmap; i < idx+1; i++ {
					cpm[field] = append(cpm[field].([]interface{}), nil)
				}
			}
			if k < len(tks.tokens)-1 {
				for i := 0; i < idx+1; i++ {
					if cpm[field].([]interface{})[i] == nil {
						cpm[field].([]interface{})[i] = map[string]interface{}{}
					}
				}
				cp = cpm[field].([]interface{})[idx]
			} else {
				//filed is last token ,set value to interface
				cpm[field].([]interface{})[idx] = value
			}
		}
	}

	return nil
}



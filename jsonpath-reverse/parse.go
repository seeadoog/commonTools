package jsonpath_reverse

import (
	"fmt"
	"regexp"
	"strconv"
)

func tokenize(exp string) ([]string, error) {
	token := make([]byte, 0, len(exp))
	tokens := make([]string, 0, 1)
	skip := false
	for i := 0; i < len(exp); i++ {
		v := exp[i]
		if v == '\\' && !skip {
			skip = true
			continue
		}
		token = append(token, v)
		if v == '.' && !skip {
			tokens = append(tokens, string(token[:len(token)-1]))
			token = token[:0]
		}
		skip = false
	}
	if len(token) > 0 {
		tokens = append(tokens, string(token))
	}

	return tokens, nil
}
var extractTokenRegexp = regexp.MustCompile(`(\w*)(\[(\d+)\])?$`)
func parseTokens(tokens []string) (tis []tokenItem, err error) {
	tis = make([]tokenItem, 0, len(tokens))
	for _, token := range tokens {
		if !extractTokenRegexp.MatchString(token){
			return nil,fmt.Errorf("invalid token:%s",token)
		}
		results:=extractTokenRegexp.FindAllStringSubmatch(token,-1)
		if len(results) >0{
			result:=results[0]
			if len(result) <4{
				return nil,fmt.Errorf("invalid token:%s",token)
			}
			ti:=tokenItem{
				key: result[1],
			}
			index:=result[3]
			if index==""{
				ti.index = -1
				tis = append(tis,ti)
				continue
			}else{
				tis = append(tis,ti)
			}

			idx,err:=strconv.Atoi(index)
			if err != nil{
				return nil,fmt.Errorf("invalid token:%s",token)
			}
			if idx <0 {
				return nil, fmt.Errorf("invalid token:%s .index must be >=0", token)
			}
			ti = tokenItem{
				index: idx,
			}

			tis = append(tis,ti)
		}
	}
	return
}

func compileTokens(s string)(tis []tokenItem,err error){
	tokens,err:=tokenize(s)
	if err != nil{
		return nil, err
	}

	return parseTokens(tokens)
}
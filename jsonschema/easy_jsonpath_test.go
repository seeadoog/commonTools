package jsonschema

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
	"testing"
)

func Test_parseTokens(t *testing.T) {
	fmt.Println(parseTokens("abc\\\\\\\\\\.efg.dfe.\\.dd"))
	fmt.Println(parseTokens(`sonf\.3`))
}

func TestReg(t *testing.T){
	var reg = regexp.MustCompile(`(\w*)(\[(\d+)\])?$`)
	str:="[44]"
	fmt.Println(reg.MatchString(str))
	res:=reg.FindAllSubmatch([]byte(str),-1)
	for i, re := range res {
		for _, s := range re {
			fmt.Println(i,string(s))
		}
	}

}

func Test_parseJpathCompiled(t *testing.T) {
	c,err:=parseJpathCompiled("bi.key")
	if err != nil{
		panicf(err.Error())
	}
	i := map[string]interface{}{
		"a":5,
		"b":"dd",
		"c":[]interface{}{1,2,3,map[string]interface{}{
			"e":5,
		},[]interface{}{3,4}},
		"bi":map[string]interface{}{

		},
	}
	fmt.Println(c.Get(i))
	fmt.Println(c.Set(i,"k"))
	fmt.Println(i)
}
func must(err error){
	if err != nil{
		panic(err)
	}
}
func TestExapmle(t *testing.T){
	b,err:=ioutil.ReadFile(`example.json`)
	if err != nil{
		panic(err)
	}
	var i Schema
	err=json.Unmarshal(b,&i)
	if err != nil{
		panic(err)
	}
	var o interface{}
	js:=`
		{
			"key1":"1",
			"status":4,
			"data":{
				"key1":"dd"
			}
		}

	`
	err = json.Unmarshal([]byte(js),&o)
	if err != nil{
		panic(err)
	}
	fmt.Println(i.Validate(o))
	fmt.Println(o)
}
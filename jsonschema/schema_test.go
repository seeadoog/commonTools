package jsonschema

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/qri-io/jsonschema"
	"testing"
)
func newType(v string)Validator{
	return Type(v)
}
func TestRaw(t *testing.T){
	f:=Prop{
		"type":Type("object"),
		"properties":Properties{
			"name":Prop{
				"type":Type("string"),
				"maxlength":MaxLength(2),
				"enums":Enums{"ets","dd"},
			},
			"age":Prop{
				"type":Type("integer"),
			},
			"sdf":Prop{

			},
			"sons":Prop{
				"type":Type("array"),
				"items":Items{
					"type":Type("string"),
				},
			},
		},
	}
	var errs = []Error{}
	i:=map[string]interface{}{
		"name":"etst",
		"age":"4",
		"fs":3,
		"sons":[]interface{}{1,2,3},
	}
	f.Validate("$",i,&errs)
	fmt.Println(errs)
}



func TestCreateNew(t *testing.T){

	var f Schema
	if err:=json.Unmarshal(schema,&f);err != nil{
		panic(err)
	}
	iv:=map[string]interface{}{
		"name":"face",
		"any":"dd",
		//"key1":"sdfsdf",
		"son":map[string]interface{}{
			"age":float64(100),
		},
		//"age":"4",
		//"fs":3,
		//"sons":[]interface{}{1,2,3},
	}
	for i:=0;i<1;i++{
		//var errs = []Error{}
		errs:=f.Validate(iv)
		fmt.Println(errs)

	}
	fmt.Println(iv)
}

func TestCreateNew2(t *testing.T){

	sc:=&jsonschema.Schema{}
	if err:=json.Unmarshal(schema,sc);err != nil{
		panic(err)
	}

	iv:=map[string]interface{}{
		"name":"face",
		"any":false,
	}
	for i:=0;i<1000000;i++{
		//var errs = []Error{}
		sc.Validate(context.Background(),iv)
		//fmt.Println(errs)
		//fmt.Println(st.Errs)
	}
}

var 	schema =[]byte(`
{
	"type":"object",
	"if":{
		"required":["key1"]
	},
	"then":{
		"required":["key2"]
	},
	"properties":{
		"name":{
			"type":"string",
            "maxLength":5,
			"minLength":3,
			"pattern":"",
			"not":{"enums":["face"]},
			"enums":["jake","face"],
			"replaceKey":"fname"
		},
		"name2":{
			"type":"string",
            "maxLength":5,
			"minLength":3,
			"enums":["jake","face"],
			"replaceKey":"fname2"
		},
		"any":{
			"anyOf":[
                 {"type":"string"},
				 {"type":"integer"}
             ]
		},
		"name3":{
			"type":"string",
            "maxLength":5,
			"minLength":3,
			"enums":["jake","face"],
			"constVal":"fname"
		},
		"son":{
			"type":"object",
			"properties":{
				"age":{
					"type":"integer",
					"maximum":100,
					"minimum":0,
					"defaultVal":15
				},
				"name":{
					"type":"string",
					"maxLength":10,
					"defaultVal":"dajj"
				}
			}
		},
		"key1":{"type":"string"},
		"key2":{"type":"integer"}
	}
}

`)

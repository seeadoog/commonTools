package jsonschema

import (
	"encoding/json"
	"fmt"
	"testing"
)
func newType(v string)Validator{
	return Type(v)
}

var schema2 = []byte(`
	
{
	"type":"object",
"properties":{
		"name":{
			"type":"string",
            "maxLength":5,
			"minLength":3,
			"pattern":"",
			"not":{"enums":["face"]},
			"enums":["jake","face"]
		
		},
		"name2":{
			"type":"string",
            "maxLength":5,
			"minLength":3,
			"enums":["jake","face"]
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
			"enums":["jake","face"]
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
		"key2":{"type":"string"}
	}
}

`)

func TestCreateNew(t *testing.T){

	var f Schema
	if err:=json.Unmarshal(schema2,&f);err != nil{
		panic(err)
	}
	iv:=map[string]interface{}{
		"name":"jake",
		"any":"dd",
		"key1":"sdfsdf",
		"key2":"dfd4s",
		"son":map[string]interface{}{
			"age":float64(100),
		},
		//"age":"4",
		//"fs":3,
		//"sons":[]interface{}{1,2,3},
	}
	type req struct {
		Name string `json:"name"`
		Any string `json:"any"`
	}
	r:=&req{
		Name: "jake2",

	}
	var errs error
	for i:=0;i<100000;i++{
		//var errs = []Error{}
		errs=f.Validate(iv)
		//errs =f.Validate(r)
		//fmt.Println(errs)

	}
	fmt.Println(r,iv,errs)

	//var a  interface{} = 1
	//var b float64 = 1
	//fmt.Println(reflect.DeepEqual(a,b))
}

//func TestCreateNew2(t *testing.T){
//
//	sc:=&jsonschema.Schema{}
//	if err:=json.Unmarshal(schema,sc);err != nil{
//		panic(err)
//	}
//
//	iv:=map[string]interface{}{
//		"name":"biaoge",
//		"any":"dd",
//		"key1":"sdfsdf",
//		"key2":"dfds",
//		"son":map[string]interface{}{
//			"age":float64(100),
//		},
//	}
//	for i:=0;i<50000;i++{
//		//var errs = []Error{}
//		sc.Validate(context.Background(),iv)
//		//fmt.Println(errs)
//		//fmt.Println(st.Errs)
//	}
//}

var 	schema =[]byte(`
{
	"type":"object",
	"if":{
		"required":["key1"]
	},
	"then":{
		"required":["key2"]
	},
	"dependencies":{
		"key1":["key2"]
	},
	"allOf":[
		{
			"if":{
				"keyMatch":{
					"name":"biaoge"
				}
			},
			"then":{
				"flexProperties":{
					"key2":{
						"type":"string",
						"maxLength":10,
						"minLength":5
					}
				}
			}
		}
	],
	"setVal":{
		"sonf\\.3name":"sname-j",
		"val1":5,
		"val2":"string",
		"val3":{
			"from":"name"
		},
		"val5":"${name}",
		"val4":{
			"from":"(append)",
			"args":[
				"hello ",
				{
					"from":"son.age"	
				}
			]
		}
	},
	"properties":{
		"name":{
			"type":"string",
            "maxLength":5,
			"minLength":3,
			"pattern":"",
			"not":{"enums":["face"]},
			"enums":["jake","face"],
			"replaceKey":"fname",
			"constVal":"gt"
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
		"key2":{"type":"string"}
	}
}

`)

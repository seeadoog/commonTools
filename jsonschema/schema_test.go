package jsonschema

import (
	//"context"
	"encoding/json"
	"fmt"
	//"github.com/qri-io/jsonschema"
	"testing"
)

var schema2 = []byte(`
	{

  "type": "object",
  "properties": {
    "a": {
      "type": "object",
      "properties": {
        "a1": {
          "type": "string",
          "maxLength": 5
        },
        "a2": {
          "type": "string",
          "maxLength": 5
        },
        "a3": {
          "type": "string",
          "maxLength": 5
        },
        "a4": {
          "type": "string"
        }
      }
    },
    "b": {
      "type": "object",
      "properties": {
        "a1": {
          "type": "string",
          "maxLength": 5
        },
        "a2": {
          "type": "string"
        },
        "a3": {
          "type": "string",
          "maxLength": 5
        },
        "a4": {
          "type": "string"
        }
      }
    },
    "c": {
      "type": "object",
      "properties": {
        "a1": {
          "type": "string",
          "maxLength": 5
        },
        "a2": {
          "type": "string"
        },
        "a3": {
          "type": "string",
          "maxLength": 5
        },
        "a4": {
          "type": "string"
        }
      }
    }
  }
}
`)

func TestCreateNew(t *testing.T) {
	ShowCompletePath = true
	var f Schema
	if err := json.Unmarshal(schema, &f); err != nil {
		panic(err)
	}
	iv := map[string]interface{}{
		"a": map[string]interface{}{
			"a1": "c",
			"a2": "1",
			"a3": "1",
			"a4": []interface{}{},
		},
		"b": map[string]interface{}{
			//"a1": "dd",
			"a2": "1",
			"a3": "1",
			"a4": float64(2),
		},
		"c": map[string]interface{}{
			"a1": "1",
			"a2": "1",
			"a3": "1",
			"a4": "5",
			"a5": float64(10),
		},
		//"age":"4",
		//"fs":3,
		//"sons":[]interface{}{1,2,3},
	}
	type req struct {
		Name string `json:"name"`
		Any  string `json:"any"`
	}
	r := &req{
		Name: "jake2",
	}
	var errs error
	for i := 0; i < 1; i++ {
		//var errs = []Error{}
		errs = f.Validate(iv)
		//errs =f.Validate(r)
		//fmt.Println(errs)

	}
	fmt.Println(r, iv, errs)

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
//	iv:=map[string]interface{}{
//		"a":map[string]interface{}{
//			"a1":"23",
//			"a2":"1",
//			"a3":"1",
//			"a4":"1",
//		},
//		"b":map[string]interface{}{
//			"a1":"1",
//			"a2":"1",
//			"a3":"1",
//			"a4":"1",
//		},
//		"c":map[string]interface{}{
//			"a1":"1",
//			"a2":"1",
//			"a3":"1",
//			"a4":"5",
//		},
//		//"age":"4",
//		//"fs":3,
//		//"sons":[]interface{}{1,2,3},
//	}
//	for i:=0;i<100000;i++{
//		//var errs = []Error{}
//		sc.Validate(context.Background(),iv)
//		//fmt.Println(errs)
//		//fmt.Println(st.Errs)
//	}
//}

var schema = []byte(`
{

  "type": "object",
  "properties": {
    "a": {
      "switch":"a1",
      "case":{
			"a":{"required":["b1","c1"]},
			"b":{"required":["b2","c2"]}
		},
		"default":{"required":["c3"]},
      "type": "object",
      "properties": {
        "a1": {
          "type": "string",
          "maxLength": 5
        },
        "a2": {
          "type": "string",
          "maxLength": 5
        },
        "a3": {
          "type": "string",
          "maxLength": 5
        },
        "a4": {"type": "string|number"}
      }
    },
    "b": {
      "type": "object",
      "if":{
			"required":["a1"]
		},
		"then":{
			"required":["b5"]
		},
		"else":{"required":["b6"]},
      "properties": {
        "a1": {
          "type": "string",
          "maxLength": 5,
          "enum":["dd"]
        },
        "a2": {
          "type": "string"
        },
        "a3": {
          "type": "string",
          "maxLength": 5
        },
        "a4": {
          "type": "string"
        }
      }
    },
    "c": {
      "type": "object",
      "properties": {
        "a1": {
          "type": "string",
          "maxLength": 0
        },
        "a2": {
          "type": "string"
        },
        "a3": {
          "type": "string",
          "maxLength": 5
        },
        "a4": {
          "type": "string"
        },
		"a5":{
			"type":"integer",
			"maximum":0
		}
      }
    }
  }
}
`)

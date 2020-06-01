## jsonschema
 
#### type  限定字段类型

取值范围：string integer number bool object array

#### properties 
当值为object 时起作用。限定object 中字段的模式，不允许出现properties 中未定义的字段

#### flexProperties

当值为object时起作用。限定object 中字段的模式，允许出现flexProperties 中未定义的字段

#### maxLength

当字段为string 或者array 类型时起作用，限定string的最大长度。（字节数）或者数组的最大长度

#### minLength 

当字段为string 或者array 类型时起作用，限定string的最小长度。（字节数）或者数组的最小长度

#### maximum 

当字段为数字类型时字作用，限定数字的最大值

#### minimum 

当字段为数字类型时起作用，限定数字的最小值

#### enum

该值类型为数组。限定值的枚举范围

````json
{
  "enum": ["1","2","3"]
}
````

#### required

该值类型为字符串数组，限定必须存在数组中声明的字段

````json
{
  "required": ["username","password"]
}
````

#### pattern 

当字段的值为字符串是起作用，pattern 的值是一个正则表达式，会校验字段是否和该正则匹配

````json
{
  "type": "array",
  "pattern": "^\\d+$"
}
````

#### items 

当字段的值为数组时起作用，用于校验数组中的每一个实体是否满足该items 中定义的模式

```json
{
  "type": "array",
  "items": {
      "type": "object",
      "properties":{
        "username": {
            "type": "string"
        }   
      }     
  }
}
```
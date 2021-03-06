package ngcfg

import (
	"fmt"
	"reflect"
	"strconv"
)

func Parse(data []byte) (*Elem, error) {
	return parse(data)
}

func Unmarshal(e *Elem, v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr {
		panic("value must be pointer")
	}
	return setObject(e, val)
}

func UnmarshalFromBytes(data []byte, v interface{}) error {
	e, err := Parse(data)
	if err != nil {
		return err
	}
	return Unmarshal(e, v)
}

func setObject(e *Elem, val reflect.Value) error {
	if val.Kind() ==reflect.Ptr{
		if val.IsNil(){
			vt:=val.Type()
			v := reflect.New(vt.Elem())
			val.Set(v)
		}
		val = val.Elem()
	}
	if _,ok:=val.Interface().(Elem);ok{
		val.Set(reflect.ValueOf(*e))
		return nil
	}
	switch val.Kind() {
	case reflect.Struct:
		t := val.Type()
		for i := 0; i < t.NumField(); i++ {
			ft := t.Field(i)
			fv := val.Field(i)
			tag := ft.Tag.Get("json")
			if tag == "" {
				tag = ft.Name
			}
			vfe := e.Get(tag)
			if vfe == nil {
				continue
			}
			switch fv.Kind() {

			case reflect.Struct, reflect.Ptr, reflect.Map:

				ele, ok := vfe.(*Elem)
				if !ok {
					return fmt.Errorf("%s is not object", tag)
				}

				if err := setObject(ele, fv); err != nil {
					return err
				}
			default:
				if fv.Kind() == reflect.Slice {
					elemKind := fv.Type().Elem().Kind()
					if elemKind == reflect.Struct || elemKind == reflect.Ptr || elemKind == reflect.Map {
						ele, ok := vfe.(*Elem)
						if !ok {
							return fmt.Errorf("%s is not object in array object", tag)
						}
						if err := setObject(ele, fv); err != nil {
							return err
						}
						break
					}
				}

				arr, ok := vfe.([]string)
				if !ok {
					return fmt.Errorf("%s is object type, want:%s", tag, fv.Kind().String())
				}
				if err := setVal(arr, fv); err != nil {
					return err
				}
			}

		}
		return nil
	case reflect.Map:

		tp := val.Type()
		if val.IsNil() {
			val.Set(reflect.MakeMap(tp))
		}
		if tp.Key().Kind() != reflect.String {
			return fmt.Errorf("key type of map must be string")
		}
		item := e.RawMap().MapItem()
		for item != nil {
			mk := item.Key
			mv := item.Val
			switch mv.(type) {
			case *Elem:
				mvType := tp.Elem()

				if mvType.Kind() == reflect.Interface {
					return fmt.Errorf("value type  interfce in map is not allowed while assgin elem  ")
				}
				var typv reflect.Value
				if mvType.Kind() == reflect.Ptr {
					typv = reflect.New(mvType.Elem())
				} else {
					typv = reflect.New(mvType)

				}
				if err := setObject(mv.(*Elem), typv); err != nil {
					return err
				}
				if mvType.Kind() == reflect.Ptr {
					val.SetMapIndex(reflect.ValueOf(mk), typv)
				} else {
					val.SetMapIndex(reflect.ValueOf(mk), typv.Elem())

				}
			case []string:
				mvType := tp.Elem()

				if mvType.Kind() == reflect.Interface {
					val.SetMapIndex(reflect.ValueOf(mk), reflect.ValueOf(mv))
					continue
				}
				mvv := reflect.New(mvType)
				if err := setVal(mv.([]string), mvv.Elem()); err != nil {
					return err
				}
				val.SetMapIndex(reflect.ValueOf(mk), mvv.Elem())

			}
			item = item.Next()
		}

		return nil
	case reflect.Slice:
		rawMap := e.RawMap()
		valType := val.Type()
		sliceElemType := valType.Elem()
		slice := reflect.MakeSlice(valType, 0, rawMap.Len())
		item := rawMap.MapItem()
		for item != nil {
			var itemVal reflect.Value
			if sliceElemType.Kind() == reflect.Ptr {
				itemVal = reflect.New(sliceElemType.Elem())
			} else {
				itemVal = reflect.New(sliceElemType)
			}
			switch item.Val.(type) {
			case []string:
				if err := setVal(item.Val.([]string), itemVal.Elem()); err != nil {
					return err
				}
			case *Elem:
				if err := setObject(item.Val.(*Elem), itemVal.Elem()); err != nil {
					return err
				}
			}
			if sliceElemType.Kind() == reflect.Ptr {
				slice = reflect.Append(slice, itemVal)
			} else {
				slice = reflect.Append(slice, itemVal.Elem())
			}

			item = item.Next()
		}
		val.Set(slice)
		return nil
	}
	return fmt.Errorf("unsupported type:%v", val.Kind().String())
}

func getSingle(val []string) string {
	target := ""
	if len(val) > 0 {
		target = val[0]
	}
	return target
}

func setVal(val []string, v reflect.Value) error {
	switch v.Kind() {
	case reflect.String:
		v.SetString(getSingle(val))
		return nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		tv := getSingle(val)
		i, err := strconv.Atoi(tv)
		if err != nil {
			return err
		}
		v.SetInt(int64(i))
		return nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint64, reflect.Uint32:
		tv := getSingle(val)
		i, err := strconv.Atoi(tv)
		if err != nil {
			return err
		}
		v.SetUint(uint64(i))
		return nil
	case reflect.Float32, reflect.Float64:
		tv := getSingle(val)
		i, err := strconv.ParseFloat(tv, 64)
		if err != nil {
			return err
		}
		v.SetFloat(i)
		return nil
	case reflect.Bool:
		tv := getSingle(val)
		bv, err := boolOf(tv)
		if err != nil {
			return err
		}
		v.SetBool(bv)
		return nil
	case reflect.Slice:
		eleType := v.Type()
		slice := reflect.MakeSlice(eleType, len(val), len(val))
		for i := 0; i < slice.Len(); i++ {
			if err := setVal([]string{val[i]}, slice.Index(i)); err != nil {
				return err
			}
		}
		v.Set(slice)
		return nil
	}
	return fmt.Errorf("unsupported type of sub value:%s", v.Kind().String())
}

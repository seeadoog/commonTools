package dbao

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
)

func MarshalFromRows(result *sql.Rows,in interface{})error{
	clms, err := result.Columns()
	if err != nil {
		return err
	}
	res := make([]interface{}, len(clms))
	for i := 0; i < len(res); i++ {
		res[i] = new(string)
	}

	mps := make([]map[string]string, 0, 5)
	for result.Next() {
		resultMap := make(map[string]string)
		if err := result.Scan(res...); err != nil {
			return err
		}
		for i, key := range clms {
			resultMap[key] = *res[i].(*string)
		}
		mps = append(mps, resultMap)
	}
	return marshal(in, mps)
}

//func newI(eleType reflect.Type,)

func marshal(in interface{}, resultMaps []map[string]string) error {
	if len(resultMaps) == 0 {
		return nil
	}
	val := reflect.ValueOf(in)
	if val.Kind() != reflect.Ptr {
		return fmt.Errorf("type of %s must be pointer", val.Type().String())
	}
	ele := val.Elem()
	valType := ele.Type()

	switch ele.Kind() {
	case reflect.Slice:
		eleType := valType.Elem()
		sv := reflect.MakeSlice(valType, 0, len(resultMaps))
		for _, resultMap := range resultMaps {
			switch eleType.Kind() {
			case reflect.Ptr:
				eleType = eleType.Elem()
				valuePtr := reflect.New(eleType)
				valuef := valuePtr.Elem()
				for i := 0; i < eleType.NumField(); i++ {
					setInto(valuef.Field(i), eleType.Field(i), resultMap)
				}
				sv = reflect.Append(sv, valuePtr)
			case reflect.Struct:
				value := reflect.New(eleType)
				valuef := value.Elem()
				for i := 0; i < eleType.NumField(); i++ {
					setInto(valuef.Field(i), eleType.Field(i), resultMap)
				}
				sv = reflect.Append(sv, value.Elem())
			}
			//sv = svreflect.AppendSlice(sv,valuePtr)
		}

		ele.Set(sv)
		//reflect.MakeSlice()
	case reflect.Struct:
		rm := resultMaps[0]
		for i := 0; i < ele.NumField(); i++ {
			fv := ele.Field(i)
			ft := valType.Field(i)

			if err := setInto(fv, ft, rm); err != nil {
				return err
			}
		}
		return nil
	default:
		return fmt.Errorf("unsupported type %s", ele.Kind().String())
	}
	return nil
}

func setInto(fv reflect.Value, ft reflect.StructField, rm map[string]string) error {
	tag := ft.Tag.Get("db")
	if tag == "" {
		tag = ft.Tag.Get("json")
	}
	if tag == "" {
		tag = ft.Name
	}
	value := rm[tag]
	if value == "" {
		return nil
	}
	switch fv.Kind() {
	case reflect.String:
		fv.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intVal, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("cannot marshall string to intval:%s->%s", tag, value)
		}
		fv.SetInt(int64(intVal))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		intVal, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("cannot marshall string to uintval:%s->%s", tag, value)
		}
		fv.SetUint(uint64(intVal))
	case reflect.Bool:
		boolv, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("cannot marshal string into bool:%s->%s", tag, value)
		}
		fv.SetBool(boolv)
	case reflect.Float32, reflect.Float64:
		flv, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("cannot marshal string into float:%s->%s", tag, value)
		}
		fv.SetFloat(flv)
	}
	return nil
}

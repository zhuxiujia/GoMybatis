package GoMybatis

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type SqlResultDecoder interface {
	Decode(s []map[string][]byte, result interface{}) error
}

type GoMybatisSqlResultDecoder struct {
	SqlResultDecoder
}

func (this GoMybatisSqlResultDecoder) Decode(s []map[string][]byte, result interface{}) error {
	if s == nil || result == nil {
		return nil
	}
	var resultV = reflect.ValueOf(result)
	var resultT = reflect.TypeOf(result)
	if resultV.Kind() == reflect.Ptr {
		resultV = resultV.Elem()
	} else {
		panic("Unmarshal only support ptr value!")
	}

	var renameMapArray = make([]map[string][]byte, 0)
	for _, v := range s {
		var m = make(map[string][]byte)
		for ik, iv := range v {
			var repleaceName = strings.ToLower(strings.Replace(ik, "_", "", -1))
			m[repleaceName] = iv
		}
		renameMapArray = append(renameMapArray, m)
	}

	var isBasicType = this.isGoBasicType(resultV.Type())

	if isBasicType {
		//single basic type
		for _, sItemMap := range s {
			if len(sItemMap) > 1 {
				return errors.New("Unmarshal one data,but sql result size find > 1 !")
			}
			for _, value := range sItemMap {
				if value == nil || len(value) == 0 {
					continue
				}
				var tItemTypeFieldTypeValue = reflect.New(resultV.Type())
				var success = this.sqlBasicTypeConvert(resultV.Type(), value, tItemTypeFieldTypeValue.Elem())
				if success {
					resultV.Set(tItemTypeFieldTypeValue.Elem())
				}
			}
		}
	} else if resultV.Kind() == reflect.Struct {
		//single struct
		if len(s) > 1 {
			return errors.New("Unmarshal one data,but sql result size find > 1 !")
		}
		for index, sItemMap := range s {
			var value = this.sqlStructConvert(resultT.Elem(), sItemMap, renameMapArray[index])
			resultV.Set(value)
		}
	} else if resultV.Kind() == reflect.Slice {
		//slice
		var resultTItemType = resultT.Elem().Elem()
		var isBasicTypeSlice = this.isGoBasicType(resultTItemType)
		if isBasicTypeSlice {
			for _, sItemMap := range s {
				for _, value := range sItemMap {
					if value == nil || len(value) == 0 {
						continue
					}
					var tItemTypeFieldTypeValue = reflect.New(resultTItemType)
					var success = this.sqlBasicTypeConvert(resultTItemType, value, tItemTypeFieldTypeValue.Elem())
					if success {
						resultV = reflect.Append(resultV, tItemTypeFieldTypeValue.Elem())
					}
				}
			}
		} else {
			for index, sItemMap := range s {
				if resultTItemType.Kind() == reflect.Struct {
					resultV = reflect.Append(resultV, this.sqlStructConvert(resultTItemType, sItemMap, renameMapArray[index]))
				}
			}
		}
	} else if resultV.Kind() == reflect.Map {
		//map
		var resultTItemType = resultT.Elem().Elem() //int,string,time.Time.....
		var isBasicTypeSlice = this.isGoBasicType(resultTItemType)
		if isBasicTypeSlice {
			if len(s) > 1 {
				return errors.New("Unmarshal one data,but sql result size find > 1 !")
			}
			for _, sItemMap := range s {
				var newResultV = reflect.MakeMap(resultT.Elem())
				for key, value := range sItemMap {
					if value == nil || len(value) == 0 {
						continue
					}
					var tItemTypeFieldTypeValue = reflect.New(resultTItemType)
					var success = this.sqlBasicTypeConvert(resultTItemType, value, tItemTypeFieldTypeValue.Elem())
					if success {
						//resultV = reflect.Append(resultV, tItemTypeFieldTypeValue.Elem())
						newResultV.SetMapIndex(reflect.ValueOf(key), tItemTypeFieldTypeValue.Elem())
					}
				}
				resultV.Set(newResultV)
			}
		} else {
			panic("[type only support map[string]interface{} and map[string]*struct{}!]")
		}
	} else {
		panic("[type only support slice and map]")
	}
	reflect.ValueOf(result).Elem().Set(resultV)

	return nil
}

func (this GoMybatisSqlResultDecoder) sqlStructConvert(resultTItemType reflect.Type, sItemMap map[string][]byte, renamedSItemMap map[string][]byte) reflect.Value {
	if resultTItemType.Kind() == reflect.Struct {
		var tItemTypeFieldTypeValue = reflect.New(resultTItemType)
		for i := 0; i < resultTItemType.NumField(); i++ {
			var tItemTypeFieldType = resultTItemType.Field(i)

			var repleaceName = tItemTypeFieldType.Name
			var value = sItemMap[repleaceName]
			if value == nil || len(value) == 0 {
				//renamed
				repleaceName = strings.ToLower(strings.Replace(tItemTypeFieldType.Name, "_", "", -1))
				value = renamedSItemMap[repleaceName]
				if value == nil || len(value) == 0 {
					continue
				}
			}
			this.sqlBasicTypeConvert(tItemTypeFieldType.Type, value, tItemTypeFieldTypeValue.Elem().Field(i))
		}
		return tItemTypeFieldTypeValue.Elem()
	} else {
		panic(resultTItemType.String() + " is not a struct obj!")
	}
}

func (this GoMybatisSqlResultDecoder) sqlBasicTypeConvert(tItemTypeFieldType reflect.Type, valueByte []byte, resultValue reflect.Value) bool {
	var value = string(valueByte)
	if tItemTypeFieldType.Kind() == reflect.String {
		resultValue.SetString(value)
	} else if tItemTypeFieldType.Kind() == reflect.Bool {
		newValue, e := strconv.ParseInt(value, 10, 64)
		if e != nil {
			return false
		}
		if newValue > 0 {
			resultValue.SetBool(true)
		} else {
			resultValue.SetBool(false)
		}
	} else if tItemTypeFieldType.Kind() == reflect.Int || tItemTypeFieldType.Kind() == reflect.Int32 || tItemTypeFieldType.Kind() == reflect.Int64 {
		newValue, e := strconv.ParseInt(value, 10, 64)
		if e != nil {
			return false
		}
		resultValue.SetInt(newValue)
	} else if tItemTypeFieldType.Kind() == reflect.Uint || tItemTypeFieldType.Kind() == reflect.Uint8 || tItemTypeFieldType.Kind() == reflect.Uint16 || tItemTypeFieldType.Kind() == reflect.Uint32 || tItemTypeFieldType.Kind() == reflect.Uint64 {
		newValue, e := strconv.ParseUint(value, 10, 64)
		if e != nil {
			return false
		}
		resultValue.SetUint(newValue)
	} else if tItemTypeFieldType.Kind() == reflect.Float32 || tItemTypeFieldType.Kind() == reflect.Float64 {
		newValue, e := strconv.ParseFloat(value, 64)
		if e != nil {
			return false
		}
		resultValue.SetFloat(newValue)
	} else if tItemTypeFieldType.String() == "time.Time" {
		newValue, e := time.Parse(string(time.RFC3339), value)
		if e != nil {
			return false
		}
		resultValue.Set(reflect.ValueOf(newValue))
	} else {
		return false
	}
	return true
}

func (this GoMybatisSqlResultDecoder) isGoBasicType(tItemTypeFieldType reflect.Type) bool {
	if tItemTypeFieldType.Kind() == reflect.Bool ||
		tItemTypeFieldType.Kind() == reflect.Int ||
		tItemTypeFieldType.Kind() == reflect.Int8 ||
		tItemTypeFieldType.Kind() == reflect.Int16 ||
		tItemTypeFieldType.Kind() == reflect.Int32 ||
		tItemTypeFieldType.Kind() == reflect.Int64 ||
		tItemTypeFieldType.Kind() == reflect.Uint ||
		tItemTypeFieldType.Kind() == reflect.Uint8 ||
		tItemTypeFieldType.Kind() == reflect.Uint16 ||
		tItemTypeFieldType.Kind() == reflect.Uint32 ||
		tItemTypeFieldType.Kind() == reflect.Uint64 ||
		tItemTypeFieldType.Kind() == reflect.Float32 ||
		tItemTypeFieldType.Kind() == reflect.Float64 ||
		tItemTypeFieldType.Kind() == reflect.String {
		return true
	}
	if tItemTypeFieldType.String() == "string" {
	} else if tItemTypeFieldType.String() == "int" {
	} else if tItemTypeFieldType.String() == "int32" {
	} else if tItemTypeFieldType.String() == "int64" {
	} else if tItemTypeFieldType.String() == "float32" {
	} else if tItemTypeFieldType.String() == "float64" {
	} else if tItemTypeFieldType.String() == "time.Time" {
	} else {
		return false
	}
	return true
}

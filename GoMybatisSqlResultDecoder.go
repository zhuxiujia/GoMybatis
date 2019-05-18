package GoMybatis

import (
	"encoding/json"
	"reflect"
	"strconv"
	"time"
)

type GoMybatisSqlResultDecoder struct {
	SqlResultDecoder
}

func (it GoMybatisSqlResultDecoder) Decode(resultMap map[string]*ResultProperty, sqlResult string, result interface{}) error {
	if sqlResult == "" || result == nil {
		return nil
	}
	var resultV = reflect.ValueOf(result)
	var resultType = resultV.Type()
	if resultV.Kind() != reflect.Ptr {
		panic("[GoMybatis] Decode only support ptr value!")
	}
	var isArray bool
	for {
		if resultType.Kind() == reflect.Ptr {
			resultType = resultType.Elem()
		} else {
			break
		}
	}
	if resultType.Kind() == reflect.Slice || resultType.Kind() == reflect.Array {
		isArray = true
	}
	if isArray {
		if it.isGoBasicType(resultType.Elem()) {
			panic("[GoMybatis]不支持的返回结果类型！，对于数组结构必须为结构体或map")
		} else {
			json.Unmarshal([]byte(sqlResult), result)
		}
	} else {
		if sqlResult != "" {
			var tempMap = []json.RawMessage{}
			json.Unmarshal([]byte(sqlResult), &tempMap)
			var resultLen = len(tempMap)
			if resultLen == 1 {
				if it.isGoBasicType(resultType) || (resultType.Kind() == reflect.Ptr && it.isGoBasicType(resultType.Elem())) {
					//basicType
					var rmap = map[string]interface{}{}
					var e = json.Unmarshal(tempMap[0], &rmap)
					if e != nil {
						return nil
					}
					for {
						if resultV.Kind() == reflect.Ptr {
							resultV = resultV.Elem()
						} else {
							break
						}
					}
					for _, v := range rmap {
						it.basicTypeConvert(resultType, v.(string), &resultV)
					}
				} else {
					json.Unmarshal(tempMap[0], result)
				}
			} else if resultLen > 1 {
				panic("[GoMybatis] sql查询结果行大于1但是返回结构体为单个")
			}
		}
	}
	return nil
}

func (it GoMybatisSqlResultDecoder) basicTypeConvert(tItemTypeFieldType reflect.Type, value string, resultValue *reflect.Value) bool {
	if value == "" {
		return false
	}
	if tItemTypeFieldType.Kind() == reflect.String {
		resultValue.SetString(value)
	} else if tItemTypeFieldType.Kind() == reflect.Bool {
		newValue, e := strconv.ParseBool(value)
		if e != nil {
			return false
		}
		resultValue.SetBool(newValue)
	} else if tItemTypeFieldType.Kind() == reflect.Int || tItemTypeFieldType.Kind() == reflect.Int32 || tItemTypeFieldType.Kind() == reflect.Int64 {
		newValue, e := strconv.ParseInt(value, 0, 64)
		if e != nil {
			return false
		}
		resultValue.SetInt(newValue)
	} else if tItemTypeFieldType.Kind() == reflect.Uint || tItemTypeFieldType.Kind() == reflect.Uint8 || tItemTypeFieldType.Kind() == reflect.Uint16 || tItemTypeFieldType.Kind() == reflect.Uint32 || tItemTypeFieldType.Kind() == reflect.Uint64 {
		newValue, e := strconv.ParseUint(value, 0, 64)
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
	} else if tItemTypeFieldType.Kind() == reflect.Struct && tItemTypeFieldType.String() == "time.Time" {
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

func (it GoMybatisSqlResultDecoder) isGoBasicType(tItemTypeFieldType reflect.Type) bool {
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
	if tItemTypeFieldType.Kind() == reflect.Struct && tItemTypeFieldType.String() == "time.Time" {
		return true
	}
	return false
}

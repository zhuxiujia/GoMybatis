package GoMybatis

import (
	"encoding/json"
	"github.com/zhuxiujia/GoMybatis/utils"
	"reflect"
	"strings"
)

type GoMybatisSqlResultDecoder struct {
	SqlResultDecoder
}

func (it GoMybatisSqlResultDecoder) Decode(resultMap map[string]*ResultProperty, sqlResult []map[string][]byte, result interface{}) error {
	if sqlResult == nil || result == nil {
		return nil
	}
	var resultV = reflect.ValueOf(result)
	if resultV.Kind() == reflect.Ptr {
		resultV = resultV.Elem()
	} else {
		panic("[GoMybatis] SqlResultDecoder only support ptr type,make sure use '*Your Type'!")
	}

	var value = []byte{}
	var sqlResultLen = len(sqlResult)
	if sqlResultLen == 0 {
		return nil
	}
	if !isArray(resultV.Kind()) {
		//single basic type
		if sqlResultLen > 1 {
			return utils.NewError("SqlResultDecoder", " Decode one result,but find database result size find > 1 !")
		}
		// base type convert
		if isBasicType(resultV.Type()) {
			for _, s := range sqlResult[0] {
				var b = strings.Builder{}
				if resultV.Kind() == reflect.String || (resultV.Kind() == reflect.Struct) {
					b.WriteString("\"")
					b.Write(s)
					b.WriteString("\"")
				} else {
					b.Write(s)
				}
				value = []byte(b.String())
				break
			}
		} else {
			var structMap, e = makeStructMap(resultV.Type())
			if e != nil {
				return e
			}
			value = makeJsonObjBytes(resultMap, sqlResult[0], structMap)
		}
	} else {
		if resultV.Type().Kind() != reflect.Array && resultV.Type().Kind() != reflect.Slice {
			return utils.NewError("SqlResultDecoder", " decode type not an struct array or slice!")
		}
		var resultVItemType = resultV.Type().Elem()
		var structMap, e = makeStructMap(resultVItemType)
		if e != nil {
			return e
		}
		var done = len(sqlResult) - 1
		var index = 0
		var jsonData = strings.Builder{}
		jsonData.WriteString("[")
		for _, v := range sqlResult {
			jsonData.Write(makeJsonObjBytes(resultMap, v, structMap))
			//write ','
			if index < done {
				jsonData.WriteString(",")
			}
			index += 1
		}
		jsonData.WriteString("]")
		value = []byte(jsonData.String())
	}
	e := json.Unmarshal(value, result)
	return e
}

func makeStructMap(itemType reflect.Type) (map[string]reflect.Type, error) {
	if itemType.Kind() != reflect.Struct {
		return nil, nil
	}
	var structMap = map[string]reflect.Type{}
	for i := 0; i < itemType.NumField(); i++ {
		var item = itemType.Field(i)
		structMap[item.Tag.Get("json")] = item.Type
	}
	return structMap, nil
}

//make an json value
func makeJsonObjBytes(resultMap map[string]*ResultProperty, sqlData map[string][]byte, structMap map[string]reflect.Type) []byte {
	var jsonData = strings.Builder{}
	jsonData.WriteString("{")

	var done = len(sqlData) - 1
	var index = 0
	for k, sqlV := range sqlData {
		jsonData.WriteString("\"")
		jsonData.WriteString(k)
		jsonData.WriteString("\":")

		var isStringType = false
		if resultMap != nil {
			var resultMapItem = resultMap[k]
			if resultMapItem != nil && (resultMapItem.LangType == "string" || resultMapItem.LangType == "time.Time") {
				isStringType = true
			}
		} else if structMap != nil {
			var v = structMap[k]
			if v != nil {
				if v.Kind() == reflect.String || v.String() == "time.Time" {
					isStringType = true
				}
			}
		} else {
			isStringType = true
		}
		if isStringType {
			jsonData.WriteString("\"")
			jsonData.WriteString(sqlVEncode(sqlV))
			jsonData.WriteString("\"")
		} else {
			if sqlV == nil || len(sqlV) == 0 {
				sqlV = []byte("null")
			}
			jsonData.Write(sqlV)
		}
		//write ','
		if index < done {
			jsonData.WriteString(",")
		}
		index += 1
	}
	jsonData.WriteString("}")
	return []byte(jsonData.String())
}

func sqlVEncode(v []byte) string {
	if v == nil || len(v) == 0 {
		return "null"
	}
	var s = string(v)
	var b, e = json.Marshal(s)
	if e != nil || len(b) == 0 {
		return "null"
	}
	s = string(b[1 : len(b)-1])
	return s
}

// is an array or slice
func isArray(kind reflect.Kind) bool {
	if kind == reflect.Slice || kind == reflect.Array {
		return true
	}
	return false
}

func isBasicType(tItemTypeFieldType reflect.Type) bool {
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

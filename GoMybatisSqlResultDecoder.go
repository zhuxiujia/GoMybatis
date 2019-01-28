package GoMybatis

import (
	"github.com/zhuxiujia/GoMybatis/utils"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type GoMybatisSqlResultDecoder struct {
	SqlResultDecoder
}

func (it GoMybatisSqlResultDecoder) Decode(resultMap map[string]*ResultProperty, sqlResult []map[string][]byte, result interface{}) error {
	if sqlResult == nil || result == nil {
		return nil
	}
	var resultV = reflect.ValueOf(result)
	var resultT = reflect.TypeOf(result)
	if resultV.Kind() == reflect.Ptr {
		resultV = resultV.Elem()
	} else {
		panic("[GoMybatis] Decode only support ptr value!")
	}
	var sqlResultLen = len(sqlResult)

	var renameMapArray = it.getRenameMapArray(sqlResult)

	if it.isGoBasicType(resultV.Type()) {
		//single basic type
		if sqlResultLen > 1 {
			return utils.NewError("SqlResultDecoder", " Decode one result,but find database result size find > 1 !")
		} else if sqlResultLen == 1 && len(sqlResult[0]) > 1 {
			return utils.NewError("SqlResultDecoder", " Decode one result,but find database result size find > 1 !")
		}
		it.convertToBasicTypeCollection(sqlResult, &resultV, resultV.Type(), resultMap)
	} else {
		switch resultV.Kind() {
		case reflect.Struct:
			//single struct
			if sqlResultLen > 1 {
				return utils.NewError("SqlResultDecoder", " Decode one result,but find database result size find > 1 !")
			}
			for index, sItemMap := range sqlResult {
				var value = it.sqlStructConvert(resultMap, resultT.Elem(), sItemMap, renameMapArray[index])
				resultV.Set(value)
			}
			break
		case reflect.Slice:
			//slice
			var resultTItemType = resultT.Elem().Elem()
			var isBasicTypeSlice = it.isGoBasicType(resultTItemType)
			if isBasicTypeSlice {
				it.convertToBasicTypeCollection(sqlResult, &resultV, resultTItemType, resultMap)
			} else {
				for index, sItemMap := range sqlResult {
					if resultTItemType.Kind() == reflect.Struct {
						resultV = reflect.Append(resultV, it.sqlStructConvert(resultMap, resultTItemType, sItemMap, renameMapArray[index]))
					} else if resultTItemType.Kind() == reflect.Map {

						var value = reflect.New(resultTItemType)
						var valueV = value.Elem()
						//map
						var resultTItemType = resultTItemType.Elem() //int,string,time.Time.....
						var isBasicTypeSlice = it.isGoBasicType(resultTItemType)
						var isInterface = resultTItemType.String() == "interface {}"
						if isBasicTypeSlice && isInterface == false {
							it.convertToBasicTypeCollection(sqlResult, &valueV, resultTItemType, resultMap)
							resultV = reflect.Append(resultV, valueV)
						} else {
							panic("[GoMybatis] Decode result type not support " + resultTItemType.String() + "!")
						}
					} else {
						panic("[GoMybatis] Decode result type not support " + resultTItemType.String() + "!")
					}
				}
			}
			break
		case reflect.Map:
			//map
			var resultTItemType = resultT.Elem().Elem() //int,string,time.Time.....
			var isBasicTypeSlice = it.isGoBasicType(resultTItemType)
			var isInterface = resultTItemType.String() == "interface {}"
			if isBasicTypeSlice && isInterface == false {
				if sqlResultLen > 1 {
					return utils.NewError("SqlResultDecoder", " Decode one result,but find database result size find > 1!")
				}
				it.convertToBasicTypeCollection(sqlResult, &resultV, resultTItemType, resultMap)
			} else {
				panic("[GoMybatis] Decode result type not support map[string]interface{}!")
			}
			break
		default:
			panic("[GoMybatis] Decode result type not support " + resultT.String() + "!")
		}
	}
	reflect.ValueOf(result).Elem().Set(resultV)
	return nil
}

func (it GoMybatisSqlResultDecoder) sqlStructConvert(resultMap map[string]*ResultProperty, resultTItemType reflect.Type, sItemMap map[string][]byte, renamedSItemMap map[string][]byte) reflect.Value {
	if resultTItemType.Kind() == reflect.Struct {
		var tItemTypeFieldTypeValue = reflect.New(resultTItemType)
		for i := 0; i < resultTItemType.NumField(); i++ {
			var tItemTypeFieldType = resultTItemType.Field(i)
			var repleaceName = tItemTypeFieldType.Name

			if !it.isGoBasicType(tItemTypeFieldType.Type) {
				//not basic type,continue
				continue
			}

			var value = sItemMap[repleaceName]
			if value == nil || len(value) == 0 {
				//renamed
				repleaceName = strings.ToLower(strings.Replace(tItemTypeFieldType.Name, "_", "", -1))
				value = renamedSItemMap[repleaceName]
				if value == nil || len(value) == 0 {
					continue
				}
			}
			var fieldValue = tItemTypeFieldTypeValue.Elem().Field(i)
			it.sqlBasicTypeConvert(repleaceName, resultMap, tItemTypeFieldType.Type, value, &fieldValue)
		}
		return tItemTypeFieldTypeValue.Elem()
	} else {
		panic(resultTItemType.String() + " is not a struct obj!")
	}
}

func (it GoMybatisSqlResultDecoder) basicTypeConvert(tItemTypeFieldType reflect.Type, valueByte []byte, resultValue *reflect.Value) bool {
	var value = string(valueByte)
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

func (it GoMybatisSqlResultDecoder) sqlBasicTypeConvert(clomnName string, resultMap map[string]*ResultProperty, tItemTypeFieldType reflect.Type, valueByte []byte, resultValue *reflect.Value) bool {
	if tItemTypeFieldType.Kind() == reflect.String {
		return it.basicTypeConvert(tItemTypeFieldType, valueByte, resultValue)
	} else if tItemTypeFieldType.Kind() == reflect.Bool {
		return it.basicTypeConvert(tItemTypeFieldType, valueByte, resultValue)
	} else if tItemTypeFieldType.Kind() == reflect.Int || tItemTypeFieldType.Kind() == reflect.Int32 || tItemTypeFieldType.Kind() == reflect.Int64 {
		return it.basicTypeConvert(tItemTypeFieldType, valueByte, resultValue)
	} else if tItemTypeFieldType.Kind() == reflect.Uint || tItemTypeFieldType.Kind() == reflect.Uint8 || tItemTypeFieldType.Kind() == reflect.Uint16 || tItemTypeFieldType.Kind() == reflect.Uint32 || tItemTypeFieldType.Kind() == reflect.Uint64 {
		return it.basicTypeConvert(tItemTypeFieldType, valueByte, resultValue)
	} else if tItemTypeFieldType.Kind() == reflect.Float32 || tItemTypeFieldType.Kind() == reflect.Float64 {
		return it.basicTypeConvert(tItemTypeFieldType, valueByte, resultValue)
	} else if tItemTypeFieldType.String() == "time.Time" {
		return it.basicTypeConvert(tItemTypeFieldType, valueByte, resultValue)
	} else {
		if resultMap != nil {
			var v = resultMap[clomnName]
			if v == nil {
				return false
			}
			if strings.EqualFold(v.Column, clomnName) || strings.EqualFold(v.Property, clomnName) {
				if v.GoType == "" {
					return false
				} else if strings.Contains(v.GoType, "string") {
					tItemTypeFieldType = StringType
				} else if strings.Contains(v.GoType, "int") {
					tItemTypeFieldType = Int64Type
				} else if strings.Contains(v.GoType, "uint") {
					tItemTypeFieldType = Uint64Type
				} else if strings.Contains(v.GoType, "time.Time") {
					tItemTypeFieldType = TimeType
				} else if strings.Contains(v.GoType, "float") {
					tItemTypeFieldType = Float64Type
				} else if strings.Contains(v.GoType, "bool") {
					tItemTypeFieldType = BoolType
				} else {
					return false
				}
				var newResultValue = reflect.New(tItemTypeFieldType)
				var newResultValueElem = newResultValue.Elem()
				if it.basicTypeConvert(tItemTypeFieldType, valueByte, &newResultValueElem) {
					resultValue.Set(newResultValue.Elem())
					return true
				} else {
					return false
				}
			}
		}
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

func (it GoMybatisSqlResultDecoder) convertToBasicTypeCollection(sourceArray []map[string][]byte, resultV *reflect.Value, itemType reflect.Type, resultMap map[string]*ResultProperty) {
	if resultV.Type().Kind() == reflect.Slice && resultV.IsValid() {
		*resultV = reflect.MakeSlice(resultV.Type(), 0, 0)
	} else if resultV.Type().Kind() == reflect.Map && resultV.IsValid() {
		*resultV = reflect.MakeMap(resultV.Type())
	} else {

	}
	for _, sItemMap := range sourceArray {
		for key, value := range sItemMap {
			if value == nil || len(value) == 0 {
				continue
			}
			var tItemTypeFieldTypeValue = reflect.New(itemType)
			var tItemTypeFieldTypeValueElem = tItemTypeFieldTypeValue.Elem()
			var success = it.sqlBasicTypeConvert(key, resultMap, itemType, value, &tItemTypeFieldTypeValueElem)
			if success {
				if resultV.Type().Kind() == reflect.Slice {
					*resultV = reflect.Append(*resultV, tItemTypeFieldTypeValue.Elem())
				} else if resultV.Type().Kind() == reflect.Map {
					resultV.SetMapIndex(reflect.ValueOf(key), tItemTypeFieldTypeValue.Elem())
				} else {
					resultV.Set(tItemTypeFieldTypeValue.Elem())
				}
			}
		}
	}
}

func (decoder GoMybatisSqlResultDecoder) getRenameMapArray(sourceArray []map[string][]byte) []map[string][]byte {
	var renameMapArray = make([]map[string][]byte, 0)
	for _, v := range sourceArray {
		var m = make(map[string][]byte)
		for ik, iv := range v {
			var repleaceName = strings.ToLower(strings.Replace(ik, "_", "", -1))
			m[repleaceName] = iv
		}
		renameMapArray = append(renameMapArray, m)
	}
	return renameMapArray
}

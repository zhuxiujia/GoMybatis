package GoMybatis

import (
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/zhuxiujia/GoMybatis/utils"
)

type GoMybatisSqlResultDecoder struct {
	SqlResultDecoder
}

func (it GoMybatisSqlResultDecoder) Decode(resultMap map[string]*ResultProperty, sqlResult QueryResult, result interface{}) error {
	if sqlResult.IsBlank() || result == nil {
		return nil
	}
	var resultV = reflect.ValueOf(result)
	var resultValue = resultV
	if resultV.Kind() == reflect.Ptr {
		resultV = resultV.Elem()
	} else {
		panic("[GoMybatis] Decode only support ptr value!")
	}
	var sqlResultLen = sqlResult.Rows()
	if it.isGoBasicType(resultV.Type()) {
		//single basic type
		if sqlResultLen > 1 {
			return utils.NewError("SqlResultDecoder", " Decode one result,but find database result size find > 1 !")
		} else if sqlResultLen == 1 && len(sqlResult.Index(0)) > 1 {
			return utils.NewError("SqlResultDecoder", " Decode one result,but find database result size find > 1 !")
		}
		it.convertToBasicTypeCollection(sqlResult.columnMap, sqlResult.Index(0), &resultV, resultMap)
	} else {
		if resultV.Kind() == reflect.Struct && sqlResultLen > 1 {
			return utils.NewError("SqlResultDecoder", " Decode one result,but find database result size find > 1 !")
		}
		for _, sItemMap := range sqlResult.data {
			it.convertToBasicTypeCollection(sqlResult.columnMap, sItemMap, &resultV, resultMap)
		}
	}
	resultValue.Elem().Set(resultV)
	return nil
}

func (it GoMybatisSqlResultDecoder) sqlStructConvert(resultMap map[string]*ResultProperty, resultTItemType reflect.Type, sItemMap map[string][]byte) reflect.Value {
	if resultTItemType.Kind() == reflect.Struct {
		var tItemTypeFieldTypeValue = reflect.New(resultTItemType)
		for i := 0; i < resultTItemType.NumField(); i++ {
			var tItemTypeFieldType = resultTItemType.Field(i)
			var jsonTag = tItemTypeFieldType.Tag.Get("json")
			var repleaceName = tItemTypeFieldType.Name

			if tItemTypeFieldType.Type.Kind() != reflect.Ptr {
				if !it.isGoBasicType(tItemTypeFieldType.Type) {
					//not basic type,continue
					continue
				}
			} else {
				if !it.isGoBasicType(tItemTypeFieldType.Type.Elem()) {
					//not basic type,continue
					continue
				}
			}
			var value = sItemMap[repleaceName]
			if value == nil || len(value) == 0 {
				//renamed
				repleaceName = jsonTag
				if repleaceName == "" {
					continue
				}
				value = sItemMap[repleaceName]
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
	} else if tItemTypeFieldType.Kind() == reflect.Int || tItemTypeFieldType.Kind() == reflect.Int8 || tItemTypeFieldType.Kind() == reflect.Int16 || tItemTypeFieldType.Kind() == reflect.Int32 || tItemTypeFieldType.Kind() == reflect.Int64 {
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
	if tItemTypeFieldType.Kind() == reflect.Ptr && valueByte != nil && len(valueByte) != 0 {
		tItemTypeFieldType = tItemTypeFieldType.Elem()
		//
		var el = resultValue.Elem()
		if el.IsValid() == false {
			resultValue.Set(reflect.New(tItemTypeFieldType))
			el = resultValue.Elem()
		}
		resultValue = &el
		return it.sqlBasicTypeConvert(clomnName, resultMap, tItemTypeFieldType, valueByte, resultValue)
	}
	if tItemTypeFieldType.Kind() == reflect.String {
		return it.basicTypeConvert(tItemTypeFieldType, valueByte, resultValue)
	} else if tItemTypeFieldType.Kind() == reflect.Bool {
		return it.basicTypeConvert(tItemTypeFieldType, valueByte, resultValue)
	} else if tItemTypeFieldType.Kind() == reflect.Int || tItemTypeFieldType.Kind() == reflect.Int8 || tItemTypeFieldType.Kind() == reflect.Int16 || tItemTypeFieldType.Kind() == reflect.Int32 || tItemTypeFieldType.Kind() == reflect.Int64 {
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
				if v.LangType == "" {
					return false
				} else if strings.Contains(v.LangType, "string") {
					tItemTypeFieldType = StringType
				} else if strings.Contains(v.LangType, "int") {
					tItemTypeFieldType = Int64Type
				} else if strings.Contains(v.LangType, "uint") {
					tItemTypeFieldType = Uint64Type
				} else if strings.Contains(v.LangType, "time.Time") {
					tItemTypeFieldType = TimeType
				} else if strings.Contains(v.LangType, "float") {
					tItemTypeFieldType = Float64Type
				} else if strings.Contains(v.LangType, "bool") {
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

//resultV:  struct,int,float,map[string]string,[]string,[]struct
func (it GoMybatisSqlResultDecoder) convertToBasicTypeCollection(columnMap map[string]SqlType, sourceMap map[string][]byte, resultV *reflect.Value, resultMap map[string]*ResultProperty) {

	var isSlice = resultV.Type().Kind() == reflect.Slice
	var isMap = resultV.Type().Kind() == reflect.Map
	var isBasicType = it.isGoBasicType(resultV.Type())
	var isStruct = resultV.Type().Kind() == reflect.Struct

	var isChildBasicType = false
	var isChildStruct = false
	var isChildInterface = false
	var isChildMap = false
	var isChildInterfaceMap = false
	if isMap || isSlice {
		var itemType = resultV.Type().Elem()
		isChildBasicType = it.isGoBasicType(itemType)
		isChildStruct = (itemType.Kind() == reflect.Struct) && !isChildBasicType
		isChildInterface = (itemType.Kind() == reflect.Interface) && !isChildBasicType
		isChildMap = !isChildBasicType && itemType.Kind() == reflect.Map
		isChildInterfaceMap = isChildMap && itemType.Elem().Kind() == reflect.Interface
	}

	if isSlice {
		//slice
		if !resultV.IsValid() || resultV.IsNil() {
			*resultV = reflect.MakeSlice(resultV.Type(), 0, 0)
		}
	} else if isMap {
		//map
		if !resultV.IsValid() || resultV.IsNil() {
			*resultV = reflect.MakeMap(resultV.Type())
		}
	} else if isBasicType {
		//basic type
	} else if isStruct {
		//struct
	} else {

	}
	var itemType = resultV.Type()
	if isBasicType {
		for key, value := range sourceMap {
			if value == nil || len(value) == 0 {
				continue
			}
			var tItemTypeFieldTypeValue = reflect.New(itemType)
			var tItemTypeFieldTypeValueElem = tItemTypeFieldTypeValue.Elem()
			var success = it.sqlBasicTypeConvert(key, resultMap, itemType, value, &tItemTypeFieldTypeValueElem)
			if success {
				resultV.Set(tItemTypeFieldTypeValue.Elem())
			}
		}
	} else if isStruct {
		var value = it.sqlStructConvert(resultMap, itemType, sourceMap)
		resultV.Set(value)
	} else if isMap {
		itemType = resultV.Type().Elem()
		if isChildBasicType {
			for key, value := range sourceMap {
				if value == nil || len(value) == 0 {
					continue
				}
				var tItemTypeFieldTypeValue = reflect.New(itemType)
				var tItemTypeFieldTypeValueElem = tItemTypeFieldTypeValue.Elem()
				var success = it.sqlBasicTypeConvert(key, resultMap, itemType, value, &tItemTypeFieldTypeValueElem)
				if success {
					resultV.SetMapIndex(reflect.ValueOf(key), tItemTypeFieldTypeValue.Elem())
				}
			}
		} else if isChildInterface {
			for key, value := range sourceMap {
				if value == nil || len(value) == 0 {
					resultV.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(struct{}{}))
					continue
				}
				point := columnMap[key].newPoint()
				if point == nil {
					continue
				}
				var tItemTypeFieldTypeValue = reflect.ValueOf(point)
				var tItemTypeFieldTypeValueElem = tItemTypeFieldTypeValue.Elem()
				var success = it.sqlBasicTypeConvert(key, resultMap, tItemTypeFieldTypeValueElem.Type(), value, &tItemTypeFieldTypeValueElem)
				if success {
					resultV.SetMapIndex(reflect.ValueOf(key), tItemTypeFieldTypeValue.Elem())
				}
			}
		} else if isChildStruct {
			panic("[GoMybatis] not supprot type struct:" + resultV.Type().String())
		} else {
			panic("[GoMybatis] not supprot type map[*]" + resultV.Type().String())
		}
	} else if isSlice {
		itemType = resultV.Type().Elem()
		if isChildBasicType {
			for key, value := range sourceMap {
				if value == nil || len(value) == 0 {
					continue
				}
				var tItemTypeFieldTypeValue = reflect.New(itemType)
				var tItemTypeFieldTypeValueElem = tItemTypeFieldTypeValue.Elem()
				var success = it.sqlBasicTypeConvert(key, resultMap, itemType, value, &tItemTypeFieldTypeValueElem)
				if success {
					*resultV = reflect.Append(*resultV, tItemTypeFieldTypeValue.Elem())
				}
			}
		} else if isChildStruct {
			var value = it.sqlStructConvert(resultMap, itemType, sourceMap)
			*resultV = reflect.Append(*resultV, value)
		} else if isChildMap {
			var mapItem = reflect.MakeMap(itemType) //todo support map[string]string -> map[string]interface{}
			for key, value := range sourceMap {
				if value == nil || len(value) == 0 {
					mapItem.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(struct{}{}))
					continue
				}
				var tItemTypeFieldTypeValue reflect.Value
				if isChildInterfaceMap {
					point := columnMap[key].newPoint()
					if point == nil {
						continue
					}
					tItemTypeFieldTypeValue = reflect.ValueOf(point)
				} else {
					tItemTypeFieldTypeValue = reflect.New(mapItem.Type().Elem())
				}
				var tItemTypeFieldTypeValueElem = tItemTypeFieldTypeValue.Elem()
				var success = it.sqlBasicTypeConvert(key, resultMap, tItemTypeFieldTypeValueElem.Type(), value, &tItemTypeFieldTypeValueElem)
				if success {
					mapItem.SetMapIndex(reflect.ValueOf(key), tItemTypeFieldTypeValueElem)
				}
			}
			*resultV = reflect.Append(*resultV, mapItem)
		} else {
			panic("[GoMybatis] not supprot type []" + itemType.String())
		}
	}
}

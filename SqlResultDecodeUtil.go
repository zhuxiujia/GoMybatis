package GoMybatis

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type SqlResultDecoder interface {
	Decode(resultMap map[string]*ResultProperty, s []map[string][]byte, result interface{}) error
}

type GoMybatisSqlResultDecoder struct {
	SqlResultDecoder
}

func (this GoMybatisSqlResultDecoder) Decode(resultMap map[string]*ResultProperty, sourceArray []map[string][]byte, result interface{}) error {
	if sourceArray == nil || result == nil {
		return nil
	}
	var resultV = reflect.ValueOf(result)
	var resultT = reflect.TypeOf(result)
	if resultV.Kind() == reflect.Ptr {
		resultV = resultV.Elem()
	} else {
		panic("[GoMybatis] Decode only support ptr value!")
	}
	var sourceArrayLen = len(sourceArray)

	var renameMapArray = this.getRenameMapArray(sourceArray)

	if this.isGoBasicType(resultV.Type()) {
		//single basic type
		if sourceArrayLen > 1 {
			return errors.New("[GoMybatis] Decode one data,but sql result size find > 1 !")
		}
		this.convertToBasicTypeCollection(sourceArray, resultV, resultV.Type(), resultMap, result)
		return nil
	}
	switch resultV.Kind() {
	case reflect.Struct:
		//single struct
		if sourceArrayLen > 1 {
			return errors.New("[GoMybatis] Decode one data,but sql result size find > 1 !")
		}
		for index, sItemMap := range sourceArray {
			var value = this.sqlStructConvert(resultMap, resultT.Elem(), sItemMap, renameMapArray[index])
			resultV.Set(value)
		}
		break
	case reflect.Slice:
		//slice
		var resultTItemType = resultT.Elem().Elem()
		var isBasicTypeSlice = this.isGoBasicType(resultTItemType)
		if isBasicTypeSlice {
			this.convertToBasicTypeCollection(sourceArray, resultV, resultTItemType, resultMap, result)
		} else {
			for index, sItemMap := range sourceArray {
				if resultTItemType.Kind() == reflect.Struct {
					resultV = reflect.Append(resultV, this.sqlStructConvert(resultMap, resultTItemType, sItemMap, renameMapArray[index]))
				}
			}
		}
		break
	case reflect.Map:
		//map
		var resultTItemType = resultT.Elem().Elem() //int,string,time.Time.....
		var isBasicTypeSlice = this.isGoBasicType(resultTItemType)
		var isInterface = resultTItemType.String() == "interface {}"
		if isBasicTypeSlice || isInterface {
			if sourceArrayLen > 1 {
				return errors.New("[GoMybatis] Decode one data,but sql result size find > 1 !")
			}
			this.convertToBasicTypeCollection(sourceArray, resultV, resultTItemType, resultMap, result)
		} else {
			panic("[GoMybatis] Decode result type only support map[string]interface{} and map[string]*struct{}!")
		}
		break
	default:
		panic("[GoMybatis] Decode result type only support slice and map")
	}
	reflect.ValueOf(result).Elem().Set(resultV)
	return nil
}

func (this GoMybatisSqlResultDecoder) sqlStructConvert(resultMap map[string]*ResultProperty, resultTItemType reflect.Type, sItemMap map[string][]byte, renamedSItemMap map[string][]byte) reflect.Value {
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
			this.sqlBasicTypeConvert(repleaceName, resultMap, tItemTypeFieldType.Type, value, tItemTypeFieldTypeValue.Elem().Field(i))
		}
		return tItemTypeFieldTypeValue.Elem()
	} else {
		panic(resultTItemType.String() + " is not a struct obj!")
	}
}

func (this GoMybatisSqlResultDecoder) basicTypeConvert(tItemTypeFieldType reflect.Type, valueByte []byte, resultValue reflect.Value) bool {
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

func (this GoMybatisSqlResultDecoder) sqlBasicTypeConvert(clomnName string, resultMap map[string]*ResultProperty, tItemTypeFieldType reflect.Type, valueByte []byte, resultValue reflect.Value) bool {
	if tItemTypeFieldType.Kind() == reflect.String {
		return this.basicTypeConvert(tItemTypeFieldType, valueByte, resultValue)
	} else if tItemTypeFieldType.Kind() == reflect.Bool {
		return this.basicTypeConvert(tItemTypeFieldType, valueByte, resultValue)
	} else if tItemTypeFieldType.Kind() == reflect.Int || tItemTypeFieldType.Kind() == reflect.Int32 || tItemTypeFieldType.Kind() == reflect.Int64 {
		return this.basicTypeConvert(tItemTypeFieldType, valueByte, resultValue)
	} else if tItemTypeFieldType.Kind() == reflect.Uint || tItemTypeFieldType.Kind() == reflect.Uint8 || tItemTypeFieldType.Kind() == reflect.Uint16 || tItemTypeFieldType.Kind() == reflect.Uint32 || tItemTypeFieldType.Kind() == reflect.Uint64 {
		return this.basicTypeConvert(tItemTypeFieldType, valueByte, resultValue)
	} else if tItemTypeFieldType.Kind() == reflect.Float32 || tItemTypeFieldType.Kind() == reflect.Float64 {
		return this.basicTypeConvert(tItemTypeFieldType, valueByte, resultValue)
	} else if tItemTypeFieldType.String() == "time.Time" {
		return this.basicTypeConvert(tItemTypeFieldType, valueByte, resultValue)
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
				if this.basicTypeConvert(tItemTypeFieldType, valueByte, newResultValue.Elem()) {
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
	if tItemTypeFieldType.Kind() == reflect.Struct && tItemTypeFieldType.String() == "time.Time" {
		return true
	}
	return false
}

func (this GoMybatisSqlResultDecoder) convertToBasicTypeCollection(sourceArray []map[string][]byte, resultV reflect.Value, itemType reflect.Type, resultMap map[string]*ResultProperty, result interface{}) {
	if resultV.Type().Kind() == reflect.Slice && resultV.IsValid() {
		resultV = reflect.MakeSlice(resultV.Type(), 0, 0)
	} else if resultV.Type().Kind() == reflect.Map && resultV.IsValid() {
		resultV = reflect.MakeMap(resultV.Type())
	} else {

	}
	for _, sItemMap := range sourceArray {
		for key, value := range sItemMap {
			if value == nil || len(value) == 0 {
				continue
			}
			var tItemTypeFieldTypeValue = reflect.New(itemType)
			var success = this.sqlBasicTypeConvert(key, resultMap, itemType, value, tItemTypeFieldTypeValue.Elem())
			if success {
				if resultV.Type().Kind() == reflect.Slice {
					resultV = reflect.Append(resultV, tItemTypeFieldTypeValue.Elem())
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

package GoMybatis

import (
	"reflect"
	"regexp"
	"strings"
)

//替换参数
func replaceArg(data string, parameters map[string]SqlArg, typeConvert SqlArgTypeConvert) string {
	if data == "" {
		return data
	}
	var defaultValue = parameters[DefaultOneArg]
	if defaultValue.Value != nil {
		var str = typeConvert.Convert(defaultValue)
		data = re.ReplaceAllString(data, str)
	}
	data = repleace(data, typeConvert, parameters)
	return data
}

var re, _ = regexp.Compile("\\#\\{[^}]*\\}")

func repleace(data string, typeConvert SqlArgTypeConvert, arg map[string]SqlArg) string {
	var findStrs = re.FindAllString(data, -1)
	var repleaceStr = ""
	for _, findStr := range findStrs {
		repleaceStr = strings.Replace(findStr, `#{`, "", -1)
		repleaceStr = strings.Replace(repleaceStr, "}", "", -1)
		if strings.Contains(repleaceStr, ",") {
			repleaceStr = strings.Split(repleaceStr, ",")[0]
		}
		data = repleaceChildFeild(data, repleaceStr, findStr, typeConvert, arg)
	}
	arg = nil
	typeConvert = nil
	return data
}

func repleaceChildFeild(data string, repleaceStr string, findStr string, typeConvert SqlArgTypeConvert, arg map[string]SqlArg) string {
	var spArr = strings.Split(repleaceStr, ".")
	if len(spArr) > 0 {
		var objStr = spArr[0]
		var fieldValue = getFeildInterface(repleaceStr, arg[objStr])
		var repleaceStr = typeConvert.Convert(fieldValue)
		data = strings.Replace(data, findStr, repleaceStr, -1)
		spArr = nil
		objStr = ""
	} else {
		repleaceStr = typeConvert.Convert(arg[repleaceStr])
		data = strings.Replace(data, findStr, repleaceStr, -1)
		return data
	}
	return data
}
func getFeildInterface(repleaceStr string, arg SqlArg) SqlArg {
	var spArr = strings.Split(repleaceStr, ".")
	if len(spArr) > 1 {
		for index, fieldName := range spArr {
			//包含子属性
			if index > 0 {
				var v = reflect.ValueOf(arg.Value).FieldByName(upperFieldFirstName(fieldName))
				arg.Value = getRealValue(v)
			}
		}
	}
	return arg
}

func getRealValue(v reflect.Value) interface{} {
	if v.Kind() == reflect.Ptr {
		if v.IsNil() == false {
			return getRealValue(v.Elem())
		}
		if v.IsNil() == false && v.CanInterface() == true {
			return v.Interface()
		} else {
			return ""
		}
	} else {
		return v.Interface()
	}
}
func upperFieldFirstName(fieldStr string) string {
	if fieldStr != "" {
		var fieldBytes = []byte(fieldStr)
		var fieldLength = len(fieldStr)
		fieldStr = strings.ToUpper(string(fieldBytes[:1])) + string(fieldBytes[1:fieldLength])
		fieldBytes = nil
	}
	return fieldStr
}

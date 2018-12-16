package GoMybatis

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

var defaultArgRegex, _ = regexp.Compile("\\$\\{[^}]*\\}")
var sqlArgRegex, _ = regexp.Compile("\\#\\{[^}]*\\}")

//替换参数
func replaceArg(data string, parameters map[string]SqlArg, typeConvert SqlArgTypeConvert) string {
	if data == "" {
		return data
	}
	var defaultValue = parameters[DefaultOneArg]
	//replace default value
	if defaultValue.Value != nil {
		var str = typeConvert.Convert(defaultValue)
		data = sqlArgRegex.ReplaceAllString(data, str)
	}
	//replace arg data
	if strings.Index(data, `#`) != -1 {
		data = replace(`#{`, sqlArgRegex, data, typeConvert, parameters)
	}
	if strings.Index(data, `$`) != -1 {
		data = replace(`${`, defaultArgRegex, data, nil, parameters)
	}
	return data
}

func replace(startChar string, regex *regexp.Regexp, data string, typeConvert SqlArgTypeConvert, arg map[string]SqlArg) string {
	var findStrs = regex.FindAllString(data, -1)
	var repleaceStr = ""
	for _, findStr := range findStrs {
		repleaceStr = strings.Replace(findStr, startChar, "", -1)
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
		if typeConvert != nil {
			var repleaceStr = typeConvert.Convert(fieldValue)
			data = strings.Replace(data, findStr, repleaceStr, -1)
		} else {
			var repleaceStr = fmt.Sprint(fieldValue.Value)
			data = strings.Replace(data, findStr, repleaceStr, -1)
		}
		spArr = nil
		objStr = ""
	} else {
		if typeConvert != nil {
			repleaceStr = typeConvert.Convert(arg[repleaceStr])
			data = strings.Replace(data, findStr, repleaceStr, -1)
		} else {
			repleaceStr = fmt.Sprint(arg[repleaceStr])
			data = strings.Replace(data, findStr, repleaceStr, -1)
		}
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

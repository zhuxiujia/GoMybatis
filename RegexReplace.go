package GoMybatis

import (
	"reflect"
	"regexp"
	"strings"
)

//替换参数
func replaceArg(data string, parameters map[string]interface{}, typeConvertFunc func(arg interface{}) string) string {
	if data == "" {
		return data
	}
	var defaultValue = parameters[DefaultOneArg]
	if defaultValue != nil {
		var str = typeConvertFunc(defaultValue)
		data = re.ReplaceAllString(data, str)
	}
	data = repleace(data, typeConvertFunc, parameters)
	return data
}

var re, _ = regexp.Compile("\\#\\{[^}]*\\}")

func repleace(data string, typeConvertFunc func(arg interface{}) string, arg map[string]interface{}) string {
	var findStrs = re.FindAllString(data, -1)
	var repleaceStr = ""
	for _, findStr := range findStrs {
		repleaceStr = strings.Replace(findStr, `#{`, "", -1)
		repleaceStr = strings.Replace(repleaceStr, "}", "", -1)
		if strings.Contains(repleaceStr, ",") {
			repleaceStr = strings.Split(repleaceStr, ",")[0]
		}
		data = repleaceChildFeild(data, repleaceStr, findStr, typeConvertFunc, arg)
	}
	arg = nil
	typeConvertFunc = nil
	return data
}

func repleaceChildFeild(data string, repleaceStr string, findStr string, typeConvertFunc func(arg interface{}) string, arg map[string]interface{}) string {
	var spArr = strings.Split(repleaceStr, ".")
	if len(spArr) > 0 {
		var objStr = spArr[0]
		var fieldValue = getFeildInterface(repleaceStr, arg[objStr])
		var repleaceStr = typeConvertFunc(fieldValue)
		data = strings.Replace(data, findStr, repleaceStr, -1)
		fieldValue = nil
		spArr = nil
		objStr = ""
	} else {
		repleaceStr = typeConvertFunc(arg[repleaceStr])
		data = strings.Replace(data, findStr, repleaceStr, -1)
		return data
	}
	return data
}
func getFeildInterface(repleaceStr string, arg interface{}) interface{} {
	var spArr = strings.Split(repleaceStr, ".")
	if len(spArr) > 1 {
		for index, fieldName := range spArr {
			//包含子属性
			if index > 0 {
				arg = reflect.ValueOf(arg).FieldByName(upperFieldFirstName(fieldName)).Interface()
			}
		}
	}
	return arg
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

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
		if strings.Contains(repleaceStr, ".") {
			var spArr = strings.Split(repleaceStr, ".")
			var objStr = spArr[0]
			var fieldStr = spArr[1]
			if fieldStr != "" {
				var fieldBytes = []byte(fieldStr)
				var fieldLength = len(fieldStr)
				fieldStr = strings.ToUpper(string(fieldBytes[:1])) + string(fieldBytes[1:fieldLength])
				fieldBytes = nil
			}
			var fieldValue = reflect.ValueOf(arg[objStr]).FieldByName(fieldStr).Interface()
			repleaceStr = typeConvertFunc(fieldValue)
			data = strings.Replace(data, findStr, repleaceStr, -1)
			fieldValue = nil
			spArr = nil
			objStr = ""
			fieldStr = ""
		} else {
			repleaceStr = typeConvertFunc(arg[repleaceStr])
			data = strings.Replace(data, findStr, repleaceStr, -1)
		}
	}
	arg = nil
	typeConvertFunc = nil
	return data
}

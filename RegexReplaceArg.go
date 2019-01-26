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
func replaceArg(data string, parameters map[string]interface{}, typeConvert SqlArgTypeConvert, engine ExpressionEngine) (string, error) {
	if data == "" {
		return data, nil
	}
	var err error
	var defaultValue = parameters[DefaultOneArg]
	//replace default value
	if defaultValue != nil {
		var str = typeConvert.Convert(SqlArg{
			Value: defaultValue,
			Type:  reflect.TypeOf(defaultValue),
		})
		data = sqlArgRegex.ReplaceAllString(data, str)
	}
	//replace arg data
	if strings.Index(data, `#`) != -1 {
		data, err = replace(`#{`, sqlArgRegex, data, typeConvert, parameters, engine)
	}
	if strings.Index(data, `$`) != -1 {
		data, err = replace(`${`, defaultArgRegex, data, nil, parameters, engine)
	}
	return data, err
}

//执行替换操作
func replace(startChar string, regex *regexp.Regexp, data string, typeConvert SqlArgTypeConvert, arg map[string]interface{}, engine ExpressionEngine) (string, error) {
	var findStrs = regex.FindAllString(data, -1)
	var repleaceStr = ""
	for _, findStr := range findStrs {
		repleaceStr = strings.Replace(findStr, startChar, "", -1)
		repleaceStr = strings.Replace(repleaceStr, "}", "", -1)
		if strings.Contains(repleaceStr, ",") {
			repleaceStr = strings.Split(repleaceStr, ",")[0]
		}
		lexer, err := engine.Lexer(repleaceStr)
		if err != nil {
			return "", err
		}
		evalData, err := engine.Eval(lexer, arg, 0)
		if err != nil {
			return "", err
		}
		if typeConvert != nil {
			repleaceStr = typeConvert.Convert(SqlArg{
				Value: evalData,
				Type:  reflect.TypeOf(evalData),
			})
		} else {
			repleaceStr = fmt.Sprint(evalData)
		}
		data = strings.Replace(data, findStr, repleaceStr, -1)
	}
	arg = nil
	typeConvert = nil
	return data, nil
}

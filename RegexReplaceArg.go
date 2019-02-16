package GoMybatis

import (
	"errors"
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
		var str = typeConvert.Convert(defaultValue, reflect.TypeOf(defaultValue))
		data = sqlArgRegex.ReplaceAllString(data, str)
	}
	//replace arg data
	if strings.Index(data, `#`) != -1 {
		data, err = replace(`#{`, FindAllExpressConvertString(data), data, typeConvert, parameters, engine)
	}
	if strings.Index(data, `$`) != -1 {
		data, err = replace(`${`, FindAllExpressString(data), data, nil, parameters, engine)
	}
	return data, err
}

//执行替换操作
func replace(startChar string, findStrs []string, data string, typeConvert SqlArgTypeConvert, arg map[string]interface{}, engine ExpressionEngine) (string, error) {
	for _, findStr := range findStrs {
		var repleaceStr = findStr
		lexer, err := engine.Lexer(repleaceStr)
		if err != nil {
			return "", errors.New(engine.Name() + ":" + err.Error())
		}
		evalData, err := engine.Eval(lexer, arg, 0)
		if err != nil {
			return "", errors.New(engine.Name() + ":" + err.Error())
		}
		if typeConvert != nil {
			repleaceStr = typeConvert.Convert(evalData, nil)
		} else {
			repleaceStr = fmt.Sprint(evalData)
		}
		data = strings.Replace(data, startChar+findStr+"}", repleaceStr, -1)
	}
	arg = nil
	typeConvert = nil
	return data, nil
}

func FindAllExpressConvertString(s string) []string {
	var finds = []string{}
	var sps = strings.Split(s, "#{")
	for _, v := range sps {
		if strings.Contains(v, "}") {
			var item = strings.Split(v, "}")[0]
			if strings.Contains(item, ",") {
				item = strings.Split(item, ",")[0]
			}
			finds = append(finds, item)
		}
	}
	return finds
}

func FindAllExpressString(s string) []string {
	var finds = []string{}
	var sps = strings.Split(s, "${")
	for _, v := range sps {
		if strings.Contains(v, "}") {
			var item = strings.Split(v, "}")[0]
			if strings.Contains(item, ",") {
				item = strings.Split(item, ",")[0]
			}
			finds = append(finds, item)
		}
	}
	return finds
}

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
	if strings.Index(data, `#`) != -1 {
		data, err = replace(`#{`, FindAllExpressConvertString(data), data, typeConvert, parameters, engine)
	}
	if strings.Index(data, `$`) != -1 {
		data, err = replace(`${`, FindAllExpressString(data), data, nil, parameters, engine)
	}
	return data, err
}

//执行替换操作
func replace(startChar string, findStrs map[string]int, data string, typeConvert SqlArgTypeConvert, arg map[string]interface{}, engine ExpressionEngine) (string, error) {
	for findStr, _ := range findStrs {
		var repleaceStr = findStr
		if strings.Contains(repleaceStr, ",") {
			repleaceStr = strings.Split(repleaceStr, ",")[0]
		}
		var evalData interface{}
		//find param arg
		var argValue = arg[findStr]
		if argValue != nil {
			evalData = argValue
		} else {
			//exec lexer
			lexer, err := engine.Lexer(repleaceStr)
			if err != nil {
				return "", errors.New(engine.Name() + ":" + err.Error())
			}
			evalData, err = engine.Eval(lexer, arg, 0)
			if err != nil {
				return "", errors.New(engine.Name() + ":" + err.Error())
			}
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

//find like #{*} value *
func FindAllExpressConvertString(str string) map[string]int {
	var finds = map[string]int{}
	var item []byte
	var lastIndex = -1
	var startIndex = -1
	var strBytes = []byte(str)
	for index, v := range strBytes {
		if v == 35 {
			lastIndex = index
		}
		if v == 123 && lastIndex == (index-1) {
			startIndex = index + 1
		}
		if v == 125 && startIndex != -1 {
			item = strBytes[startIndex:index]
			finds[string(item)] = 1
			item = nil
			startIndex = -1
			lastIndex = -1
		}
	}
	item = nil
	strBytes = nil
	return finds
}

//find like ${*} value *
func FindAllExpressString(str string) map[string]int {
	var finds = map[string]int{}
	var item []byte
	var lastIndex = -1
	var startIndex = -1
	var strBytes = []byte(str)
	for index, v := range str {
		if v == 36 {
			lastIndex = index
		}
		if v == 123 && lastIndex == (index-1) {
			startIndex = index + 1
		}
		if v == 125 && startIndex != -1 {
			item = strBytes[startIndex:index]
			finds[string(item)] = 1
			item = nil
			startIndex = -1
			lastIndex = -1
		}
	}
	item = nil
	strBytes = nil
	return finds
}

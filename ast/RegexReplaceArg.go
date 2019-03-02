package ast

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var defaultArgRegex, _ = regexp.Compile("\\$\\{[^}]*\\}")
var sqlArgRegex, _ = regexp.Compile("\\#\\{[^}]*\\}")

//执行替换操作
func Replace(startChar string, findStrs []string, data string, typeConvert SqlArgTypeConvert, arg map[string]interface{}, engine ExpressionEngine) (string, error) {
	for _, findStr := range findStrs {
		var evalData interface{}
		//find param arg
		var argValue = arg[findStr]
		if argValue != nil {
			evalData = argValue
		} else {
			//exec lexer
			var err error
			evalData, err = engine.LexerAndEval(findStr, arg)
			if err != nil {
				return "", errors.New(engine.Name() + ":" + err.Error())
			}
		}
		var resultStr string
		if typeConvert != nil {
			resultStr = typeConvert.Convert(evalData, nil)
		} else {
			resultStr = fmt.Sprint(evalData)
		}
		data = strings.Replace(data, startChar+findStr+"}", resultStr, -1)
	}
	arg = nil
	typeConvert = nil
	return data, nil
}

//find like #{*} value *
func FindAllExpressConvertString(str string) []string {
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

			//去掉逗号之后的部分
			if bytes.Contains(item, []byte(",")) {
				item = bytes.Split(item, []byte(","))[0]
			}

			finds[string(item)] = 1
			item = nil
			startIndex = -1
			lastIndex = -1
		}
	}
	item = nil
	strBytes = nil

	var strs = []string{}
	for k, _ := range finds {
		strs = append(strs, k)
	}
	return strs
}

//find like ${*} value *
func FindAllExpressString(str string) []string {
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

			//去掉逗号之后的部分
			if bytes.Contains(item, []byte(",")) {
				item = bytes.Split(item, []byte(","))[0]
			}

			finds[string(item)] = 1
			item = nil
			startIndex = -1
			lastIndex = -1
		}
	}
	item = nil
	strBytes = nil

	var strs = []string{}
	for k, _ := range finds {
		strs = append(strs, k)
	}
	return strs
}

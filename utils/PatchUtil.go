package utils

import (
	"regexp"
	"strings"
)
//补丁：临时解决test表达式，包含<符号时etree 无法解析的问题
func FixTestExpressionSymbol(bytes *[]byte) {
	var byteStr = string(*bytes)
	var findStrs = getTestRegex().FindAllString(byteStr, -1)
	for _, findStr := range findStrs {
		var newStr = string(findStr)
		newStr = strings.Replace(newStr, "<", "&lt;", -1)
		byteStr = strings.Replace(byteStr, findStr, newStr, -1)
	}
	*bytes = []byte(byteStr)
}

var testRegex *regexp.Regexp
func getTestRegex() *regexp.Regexp {
	if testRegex == nil {
		testRegex, _ = regexp.Compile(`test=".*<.*"`)
	}
	return testRegex
}

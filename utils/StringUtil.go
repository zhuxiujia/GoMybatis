package utils

import (
	"fmt"
	"strings"
)

//首字母转大写
func UpperFieldFirstName(fieldStr string) string {
	if fieldStr != "" {
		var fieldBytes = []byte(fieldStr)
		var fieldLength = len(fieldStr)
		fieldStr = strings.ToUpper(string(fieldBytes[:1])) + string(fieldBytes[1:fieldLength])
		fieldBytes = nil
	}
	return fieldStr
}

//首字母转小写
func LowerFieldFirstName(fieldStr string) string {
	if fieldStr != "" {
		var fieldBytes = []byte(fieldStr)
		var fieldLength = len(fieldStr)
		fieldStr = strings.ToLower(string(fieldBytes[:1])) + string(fieldBytes[1:fieldLength])
		fieldBytes = nil
	}
	return fieldStr
}

// format array [1,2,3,""] to '[1,2,3,]'
func SprintArray(array_or_slice []interface{}) string {
	return strings.Replace(strings.Trim(fmt.Sprint(array_or_slice), ""), " ", ",", -1)
}

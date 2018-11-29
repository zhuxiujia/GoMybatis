package GoMybatis

import (
	"bytes"
	"fmt"
	"reflect"
	"time"
)

const Adapter_DateType = `time.Time`
const Adapter_FormateDate = `2006-01-02 15:04:05`

//表达式类型(基本类型)转换函数
var DefaultExpressionTypeConvertFunc = func(arg interface{}) interface{} {
	var t = reflect.TypeOf(arg)
	if t.Kind() == reflect.Struct && t.String() == Adapter_DateType {
		return arg.(time.Time).Nanosecond()
	}
	return arg
}

//默认sql参数(基本类型)转换函数
var DefaultSqlTypeConvertFunc = func(arg interface{}) string {
	if arg == nil {
		return ""
	}
	var t = reflect.TypeOf(arg)
	if t.Kind() == reflect.Struct && t.String() == Adapter_DateType {
		arg = arg.(time.Time).Format(Adapter_FormateDate)
	}
	if t.Kind() == reflect.Bool {
		if arg.(bool) {
			arg = 1
		} else {
			arg = 0
		}
	}
	if t.Kind() == reflect.String {
		var argStr bytes.Buffer
		argStr.WriteString(`'`)
		argStr.WriteString(toString(arg))
		argStr.WriteString(`'`)
		return argStr.String()
	}
	return toString(arg)
}

func toString(value interface{}) string {
	if value == nil {
		return ""
	}
	return fmt.Sprint(value)
}

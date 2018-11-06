package GoMybatis

import (
	"reflect"
	"time"
	"strconv"
)

var DefaultExpressionTypeConvertFunc = func(arg interface{}) interface{} {
	if reflect.TypeOf(arg).String() == `time.Time` {
		return arg.(time.Time).Nanosecond()
	}
	return arg
}

var DefaultSqlTypeConvertFunc = func(arg interface{}) string {
	var t=reflect.TypeOf(arg)
	if t.String() == `time.Time` {
		arg = arg.(time.Time).Format(`2006-01-02 15:04:05`)
	}
	if t.String()  == `time.Time`|| t.String()==`string`{
		return `'`+toString(arg)+`'`
	}
	return toString(arg)
}


func toString(value interface{}) string {
	if value == nil {
		return ""
	}
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Int {
		string := strconv.Itoa(value.(int))
		return string
	} else if v.Kind() == reflect.Int64 {
		string := strconv.FormatInt(value.(int64), 10)
		return string
	} else if v.Kind() == reflect.Float32 {
		string := strconv.FormatFloat(value.(float64), 'f', 8, 64)
		return string
	} else if v.Kind() == reflect.Float64 {
		string := strconv.FormatFloat(value.(float64), 'f', 8, 64)
		return string
	} else if v.Kind() == reflect.String {
		return value.(string)
	} else {
		return ""
	}
}
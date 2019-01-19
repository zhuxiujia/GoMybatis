package GoMybatis

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

const Adapter_DateType = `time.Time`
const Adapter_FormateDate = `2006-01-02 15:04:05`

//表达式类型(基本类型)转换函数
type ExpressionTypeConvert interface {
	Convert(arg SqlArg) interface{}
}

//表达式类型(基本类型)转换函数
type SqlArgTypeConvert interface {
	Convert(arg SqlArg) string
}

//表达式类型转换器
type GoMybatisExpressionTypeConvert struct {
}

//表达式类型转换器
func (this GoMybatisExpressionTypeConvert) Convert(arg SqlArg) interface{} {
	if arg.Type.Kind() == reflect.Struct && arg.Type.String() == Adapter_DateType {
		return arg.Value.(time.Time).Nanosecond()
	}
	return arg.Value
}

//Sql内容类型转换器
type GoMybatisSqlArgTypeConvert struct {
}

//Sql内容类型转换器
func (this GoMybatisSqlArgTypeConvert) Convert(arg SqlArg) string {
	var argValue = arg.Value
	var argType = arg.Type
	if argValue == nil {
		return "''"
	}
	switch argType.Kind() {
	case reflect.Bool:
		if argValue.(bool) {
			argValue = strconv.FormatBool(true)
		} else {
			argValue = strconv.FormatBool(false)
		}
		break
	case reflect.String:
		var argStr bytes.Buffer
		argStr.WriteString(`'`)
		argStr.WriteString(this.toString(&arg))
		argStr.WriteString(`'`)
		return argStr.String()
	case reflect.Struct:
		if argType.String() == Adapter_DateType {
			var argStr bytes.Buffer
			argStr.WriteString(`'`)
			argStr.WriteString(argValue.(time.Time).Format(Adapter_FormateDate))
			argStr.WriteString(`'`)
			return argStr.String()
		}
		break
	}
	return this.toString(&arg)
}

func (this GoMybatisSqlArgTypeConvert) toString(value *SqlArg) string {
	if value.Value == nil {
		return ""
	}
	return fmt.Sprint(value.Value)
}

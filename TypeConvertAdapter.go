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

type GoMybatisExpressionTypeConvert struct {
	ExpressionTypeConvert
}

func (this GoMybatisExpressionTypeConvert) Convert(arg SqlArg) interface{} {
	if arg.Type.Kind() == reflect.Struct && arg.Type.String() == Adapter_DateType {
		return arg.Value.(time.Time).Nanosecond()
	}
	return arg.Value
}

type GoMybatisSqlArgTypeConvert struct {
	SqlArgTypeConvert
}

func (this GoMybatisSqlArgTypeConvert) Convert(arg SqlArg) string {
	if arg.Value == nil {
		return ""
	}
	if arg.Type.Kind() == reflect.Struct && arg.Type.String() == Adapter_DateType {
		arg.Value = arg.Value.(time.Time).Format(Adapter_FormateDate)
	}
	if arg.Type.Kind() == reflect.Bool {
		if arg.Value.(bool) {
			arg.Value = strconv.FormatBool(true)
		} else {
			arg.Value = strconv.FormatBool(false)
		}
	}
	if arg.Type.Kind() == reflect.String {
		var argStr bytes.Buffer
		argStr.WriteString(`'`)
		argStr.WriteString(this.toString(&arg))
		argStr.WriteString(`'`)
		return argStr.String()
	}
	return this.toString(&arg)
}

func (this GoMybatisSqlArgTypeConvert) toString(value *SqlArg) string {
	if value.Value == nil {
		return ""
	}
	return fmt.Sprint(value.Value)
}

package GoMybatis

import (
	"reflect"
	"time"
)

const Adapter_DateType = `time.Time`

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

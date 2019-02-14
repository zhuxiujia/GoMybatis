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
func (it GoMybatisExpressionTypeConvert) Convert(arg interface{},argType reflect.Type) interface{} {
	if argType==nil{
		argType=reflect.TypeOf(arg)
	}
	if argType.Kind() == reflect.Struct && argType.String() == Adapter_DateType {
		return arg.(time.Time).Nanosecond()
	}
	return arg
}

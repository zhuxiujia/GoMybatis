package GoMybatis

import "reflect"

//表达式类型(基本类型)转换函数
type ExpressionTypeConvert interface {
	Convert(arg interface{},argType reflect.Type) interface{}
}

package ast

//表达式类型(基本类型)转换函数
type SqlArgTypeConvert interface {
	Convert(arg interface{}) string
}

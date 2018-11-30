package GoMybatis

import "reflect"

type SqlArg struct {
	Value interface{}
	Type  reflect.Type
}

package GoMybatis

import (
	"reflect"
)

type ReturnType struct {
	ErrorType     *reflect.Type
	ReturnOutType *reflect.Type
	ReturnIndex   int //返回数据位置索引
	NumOut        int //返回总数
}

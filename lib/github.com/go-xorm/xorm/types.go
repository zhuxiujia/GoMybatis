package xorm

import (
	"reflect"

	"github.com/zhuxiujia/GoMybatis/lib/github.com/go-xorm/core"
)

var (
	ptrPkType = reflect.TypeOf(&core.PK{})
	pkType    = reflect.TypeOf(core.PK{})
)

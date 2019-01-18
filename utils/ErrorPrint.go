package utils

import (
	"errors"
	"fmt"
)

const PackageName  = "com.github.zhuxiujia.GoMybatis"

func NewError(StructName string,args ...interface{}) error {
	return errors.New(fmt.Sprint("[GoMybatis] ", PackageName, ".",StructName,": ", args))
}

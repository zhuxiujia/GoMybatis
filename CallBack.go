package GoMybatis

import "reflect"

//sql 运行时 回调
type CallBack struct {
	//sql 查询之前执行，可以写指针改变sql内容(func 可以为nil)
	BeforeQuery func(args []reflect.Value, sqlorArgs *string)
	//sql 查询之前执行，可以写指针改变sql内容(func 可以为nil)
	BeforeExec func(args []reflect.Value, sqlorArgs *string)

	//sql 查询之后执行，可以写指针改变返回结果(func 可以为nil)
	AfterQuery func(args []reflect.Value, sqlorArgs string, result *[]map[string][]byte, err *error)
	//sql 查询之后执行，可以写指针改变返回结果(func 可以为nil)
	AfterExec func(args []reflect.Value, sqlorArgs string, result *Result, err *error)
}

package GoMybatis

//sql 运行时 回调
type CallBack struct {
	//sql 查询之前执行，可以写指针改变sql内容(func 可以为nil)
	BeforeQuery func(sqlorArgs *string)
	//sql 查询之前执行，可以写指针改变sql内容(func 可以为nil)
	BeforeExec func(sqlorArgs *string)

	//sql 查询之后执行，可以写指针改变返回结果(func 可以为nil)
	AfterQuery func(sqlorArgs string, result *[]map[string][]byte, err *error)
	//sql 查询之后执行，可以写指针改变返回结果(func 可以为nil)
	AfterExec func(sqlorArgs string, result *Result, err *error)
}

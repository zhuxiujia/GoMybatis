# 表达式执行引擎包

## 引擎必须实现AST包 github.com/zhuxiujia/GoMybatis/ast 下  ast.ExpressionEngine
```
type ExpressionEngine interface {
	//引擎名称
	Name() string
	//编译一个表达式
	//参数：lexerArg 表达式内容
	//返回：interface{} 编译结果,error 错误
	Lexer(lexerArg string) (interface{}, error)

	//执行一个表达式
	//参数：lexerResult=编译结果，arg=参数
	//返回：执行结果，错误
	Eval(lexerResult interface{}, arg interface{}, operation int) (interface{}, error)

    
	LexerAndEval(lexerArg string,arg interface{})  (interface{}, error)
}
```
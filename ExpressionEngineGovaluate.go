package GoMybatis

import "github.com/zhuxiujia/GoMybatis/lib/github.com/Knetic/govaluate"

type ExpressionEngineGovaluate struct {
}

func (this *ExpressionEngineGovaluate) Name() string {
	return "ExpressionEngineGovaluate"
}
//编译一个表达式
//参数：lexerArg 表达式内容
//返回：interface{} 编译结果,error 错误
func (this *ExpressionEngineGovaluate) Lexer(expression string) (interface{}, error) {
	return govaluate.NewEvaluableExpression(expression)
}

//执行一个表达式
//参数：lexerResult=编译结果，arg=参数
//返回：执行结果，错误
func (this *ExpressionEngineGovaluate) Eval(compileResult interface{}, arg interface{}, operation int) (interface{}, error) {
	return compileResult.(*govaluate.EvaluableExpression).Evaluate(arg.(map[string]interface{}))
}

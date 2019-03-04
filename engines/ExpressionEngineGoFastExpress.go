package engines

import (
	"github.com/zhuxiujia/GoMybatis/lib/github.com/zhuxiujia/GoFastExpress"
	"strings"
)

type ExpressionEngineGoExpress struct {
}

//引擎名称
func (it *ExpressionEngineGoExpress) Name() string {
	return "ExpressionEngineGoExpress"
}

//编译一个表达式
//参数：lexerArg 表达式内容
//返回：interface{} 编译结果,error 错误
func (it *ExpressionEngineGoExpress) Lexer(expression string) (interface{}, error) {
	expression = it.repleaceExpression(expression)
	var result, err = GoFastExpress.Parser(expression)
	return result, err
}

//执行一个表达式
//参数：lexerResult=编译结果，arg=参数
//返回：执行结果，错误
func (it *ExpressionEngineGoExpress) Eval(lexerResult interface{}, arg interface{}, operation int) (interface{}, error) {
	output, err := lexerResult.(GoFastExpress.Node).Eval(arg)
	return output, err
}

func (it *ExpressionEngineGoExpress) LexerAndEval(lexerArg string, arg interface{}) (interface{}, error) {
	var lex, err = it.Lexer(lexerArg)
	if err != nil {
		return nil, err
	}
	return it.Eval(lex, arg, 0)
}

//替换表达式中的值 and,or,参数 替换为实际值
func (it *ExpressionEngineGoExpress) repleaceExpression(expression string) string {
	if expression == "" {
		return expression
	}
	expression = strings.Replace(expression, ` and `, " && ", -1)
	expression = strings.Replace(expression, ` or `, " || ", -1)
	return expression
}

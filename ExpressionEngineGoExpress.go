package GoMybatis

import (
	"github.com/zhuxiujia/GoMybatis/lib/github.com/GoExpress"
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
	var result, err = GoExpress.Parser(expression)
	return result, err
}

//执行一个表达式
//参数：lexerResult=编译结果，arg=参数
//返回：执行结果，错误
func (it *ExpressionEngineGoExpress) Eval(lexerResult interface{}, arg interface{}, operation int) (interface{}, error) {
	output, err := lexerResult.(GoExpress.Node).Eval(arg)
	return output, err
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

func (it *ExpressionEngineGoExpress) split(str string) (stringItems []string) {
	if str == "" {
		return nil
	}
	var andStrings = strings.Split(str, " && ")
	if andStrings == nil {
		return nil
	}
	var newStrings []string
	for _, v := range andStrings {
		var orStrings = strings.Split(v, " || ")
		if orStrings == nil {
			continue
		}
		for _, orStr := range orStrings {
			if newStrings == nil {
				newStrings = make([]string, 0)
			}
			if orStr == "" {
				continue
			}
			newStrings = append(newStrings, orStr)
		}
	}
	return newStrings
}

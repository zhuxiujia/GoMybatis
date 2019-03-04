package engines

import (
	"github.com/zhuxiujia/GoMybatis/lib/github.com/Knetic/govaluate"
	"reflect"
	"strings"
)

type ExpressionEngineGovaluate struct {
}

func (it *ExpressionEngineGovaluate) Name() string {
	return "ExpressionEngineGovaluate"
}

//编译一个表达式
//参数：lexerArg 表达式内容
//返回：interface{} 编译结果,error 错误
func (it *ExpressionEngineGovaluate) Lexer(expression string) (interface{}, error) {
	expression = it.repleaceExpression(expression)
	return govaluate.NewEvaluableExpression(expression)
}

//执行一个表达式
//参数：lexerResult=编译结果，arg=参数
//返回：执行结果，错误
func (it *ExpressionEngineGovaluate) Eval(compileResult interface{}, arg interface{}, operation int) (interface{}, error) {
	var argMap = arg.(map[string]interface{})
	for k, v := range argMap {
		if v != nil {
			var reflectV = reflect.ValueOf(v)
			if reflectV.IsValid() == false || (reflectV.Kind() == reflect.Ptr && reflectV.IsNil()) {
				argMap[k] = nil
			}
		}
	}
	argMap["nil"] = nil
	argMap["null"] = nil
	return compileResult.(*govaluate.EvaluableExpression).Evaluate(argMap)
}

func (it *ExpressionEngineGovaluate) LexerAndEval(lexerArg string, arg interface{}) (interface{}, error) {
	var lex, err = it.Lexer(lexerArg)
	if err != nil {
		return nil, err
	}
	return it.Eval(lex, arg, 0)
}

//替换表达式中的值 and,or,参数 替换为实际值
func (it *ExpressionEngineGovaluate) repleaceExpression(expression string) string {
	if expression == "" {
		return expression
	}
	expression = strings.Replace(expression, ` and `, " && ", -1)
	expression = strings.Replace(expression, ` or `, " || ", -1)
	return expression
}

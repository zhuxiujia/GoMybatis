package GoMybatis

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
	var argMap=arg.(map[string]interface{})
	for k,v:=range argMap{
		if v!=nil{
			var reflectV=reflect.ValueOf(v)
			if reflectV.IsValid()==false || reflectV.IsNil(){
				argMap[k]=nil
			}
		}
	}
	argMap["nil"]=nil
	argMap["null"]=nil
	return compileResult.(*govaluate.EvaluableExpression).Evaluate(argMap)
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

func (it *ExpressionEngineGovaluate) split(str string) (stringItems []string) {
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

//Lexer缓存,可不提供。
func (it *ExpressionEngineGovaluate) LexerCache() ExpressionEngineLexerCache {
	return nil
}

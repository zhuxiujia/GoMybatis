package GoMybatis

import (
	"github.com/zhuxiujia/GoMybatis/utils"
)

type ExpressionEngineProxy struct {
	expressionEngineLexerCache ExpressionEngineLexerCache //lexer缓存接口，默认使用ExpressionEngineLexerMapCache
	expressionEngine           ExpressionEngine
	lexerCacheable             bool //是否使用lexer缓存,默认false
}

//engine ：表达式引擎,useLexerCache：是否缓存Lexer表达式编译结果
func (ExpressionEngineProxy) New(engine ExpressionEngine, useLexerCache bool) ExpressionEngineProxy {
	var it = ExpressionEngineProxy{
		expressionEngine: engine,
		lexerCacheable:   useLexerCache,
	}
	if it.expressionEngineLexerCache == nil {
		var cache = ExpressionEngineLexerMapCache{}.New()
		it.SetLexerCache(&cache)
	}
	return it
}

//引擎名称
func (it *ExpressionEngineProxy) SetExpressionEngine(engine ExpressionEngine) {
	it.expressionEngine = engine
}

//引擎名称
func (it ExpressionEngineProxy) Name() string {
	if it.expressionEngine == nil {
		return ""
	}
	return it.expressionEngine.Name()
}

//编译一个表达式
//参数：lexerArg 表达式内容
//返回：interface{} 编译结果,error 错误
func (it *ExpressionEngineProxy) Lexer(expression string) (interface{}, error) {
	if it.expressionEngine == nil {
		return nil, utils.NewError("ExpressionEngineProxy", "ExpressionEngineProxy not init for ExpressionEngineProxy{}.New(...)")
	}

	if it.lexerCacheable {
		if it.LexerCache() == nil {
			panic(utils.NewError("ExpressionEngineProxy", "lexerCacheable =true! lexerCache can not be nil! you must set the cache!"))
		}
		//如果 提供缓存，则使用缓存
		cacheResult, cacheErr := it.LexerCache().Get(expression)
		if cacheErr != nil {
			return nil, cacheErr
		}
		if cacheResult != nil {
			return cacheResult, nil
		}
	}
	var result, err = it.expressionEngine.Lexer(expression)
	if it.lexerCacheable && err == nil {
		//如果 提供缓存，则使用缓存
		it.LexerCache().Set(expression, result)
	}
	return result, err
}

//执行一个表达式
//参数：lexerResult=编译结果，arg=参数
//返回：执行结果，错误
func (it *ExpressionEngineProxy) Eval(lexerResult interface{}, arg interface{}, operation int) (interface{}, error) {
	if it.expressionEngine == nil {
		return nil, utils.NewError("ExpressionEngineProxy", "ExpressionEngineProxy not init for ExpressionEngineProxy{}.New(...)")
	}
	return it.expressionEngine.Eval(lexerResult, arg, operation)
}

func (it *ExpressionEngineProxy) LexerCache() ExpressionEngineLexerCache {
	return it.expressionEngineLexerCache
}

func (it *ExpressionEngineProxy) SetLexerCache(cache ExpressionEngineLexerCache) {
	it.expressionEngineLexerCache = cache
}

func (it *ExpressionEngineProxy) SetUseLexerCache(isUseCache bool) error {
	it.lexerCacheable = isUseCache
	return nil
}
func (it *ExpressionEngineProxy) LexerCacheable() bool {
	return it.lexerCacheable
}

//执行
func (it *ExpressionEngineProxy) LexerAndEval(expression string, arg map[string]interface{}) (interface{}, error) {

	var funcItem = arg["func_"+expression]
	if funcItem != nil {
		var f = funcItem.(func(arg map[string]interface{}) interface{})
		return f(arg), nil
	}
	ifElementevalExpression, err := it.Lexer(expression)
	if err != nil {
		return false, utils.NewError("ExpressionEngineProxy", err)
	}
	result, err := it.Eval(ifElementevalExpression, arg, 0)
	if err != nil {
		return false, utils.NewError("ExpressionEngineProxy", err)
	}
	return result, nil
}

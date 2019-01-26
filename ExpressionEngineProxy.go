package GoMybatis

import "github.com/zhuxiujia/GoMybatis/utils"

type ExpressionEngineProxy struct {
	expressionEngine ExpressionEngine
	lexerCacheable   bool //是否使用lexer缓存,默认false
}

//engine ：表达式引擎,useLexerCache：是否缓存Lexer表达式编译结果
func (ExpressionEngineProxy) New(engine ExpressionEngine, useLexerCache bool) ExpressionEngineProxy {
	return ExpressionEngineProxy{
		expressionEngine: engine,
		lexerCacheable:   useLexerCache,
	}
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
	if it.expressionEngine.LexerCache() != nil && it.lexerCacheable {
		//如果 提供缓存，则使用缓存
		cacheResult, cacheErr := it.expressionEngine.LexerCache().Get(expression)
		if cacheErr != nil {
			return nil, cacheErr
		}
		if cacheResult != nil {
			return cacheResult, nil
		}
	}
	var result, err = it.expressionEngine.Lexer(expression)
	if it.expressionEngine.LexerCache() != nil && it.lexerCacheable {
		//如果 提供缓存，则使用缓存
		it.expressionEngine.LexerCache().Set(expression, result)
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
	if it.expressionEngine == nil {
		return nil
	}
	return it.expressionEngine.LexerCache()
}

func (it *ExpressionEngineProxy) SetUseLexerCache(isUseCache bool) error {
	it.lexerCacheable = isUseCache
	return nil
}
func (it *ExpressionEngineProxy) LexerCacheable() bool {
	return it.lexerCacheable
}

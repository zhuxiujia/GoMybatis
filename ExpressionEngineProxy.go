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
func (this ExpressionEngineProxy) Name() string {
	if this.expressionEngine == nil {
		return ""
	}
	return this.expressionEngine.Name()
}

//编译一个表达式
//参数：lexerArg 表达式内容
//返回：interface{} 编译结果,error 错误
func (this *ExpressionEngineProxy) Lexer(expression string) (interface{}, error) {
	if this.expressionEngine == nil {
		return nil, utils.NewError("ExpressionEngineProxy", "ExpressionEngineProxy not init for ExpressionEngineProxy{}.New(...)")
	}
	if this.expressionEngine.LexerCache() != nil && this.lexerCacheable {
		//如果 提供缓存，则使用缓存
		cacheResult, cacheErr := this.expressionEngine.LexerCache().Get(expression)
		if cacheErr != nil {
			return nil, cacheErr
		}
		if cacheResult != nil {
			return cacheResult, nil
		}
	}
	var result, err = this.expressionEngine.Lexer(expression)
	if this.expressionEngine.LexerCache() != nil && this.lexerCacheable {
		//如果 提供缓存，则使用缓存
		this.expressionEngine.LexerCache().Set(expression, result)
	}
	return result, err
}

//执行一个表达式
//参数：lexerResult=编译结果，arg=参数
//返回：执行结果，错误
func (this *ExpressionEngineProxy) Eval(lexerResult interface{}, arg interface{}, operation int) (interface{}, error) {
	if this.expressionEngine == nil {
		return nil, utils.NewError("ExpressionEngineProxy", "ExpressionEngineProxy not init for ExpressionEngineProxy{}.New(...)")
	}
	return this.expressionEngine.Eval(lexerResult, arg, operation)
}

func (this *ExpressionEngineProxy) LexerCache() ExpressionEngineLexerCache {
	if this.expressionEngine == nil {
		return nil
	}
	return this.expressionEngine.LexerCache()
}

func (this *ExpressionEngineProxy) SetUseLexerCache(isUseCache bool) error {
	this.lexerCacheable = isUseCache
	return nil
}
func (this *ExpressionEngineProxy) LexerCacheable() bool {
	return this.lexerCacheable
}

package GoMybatis

import (
	"github.com/zhuxiujia/GoMybatis/utils"
)

type ExpressionEngineLexerMapCache struct {
	mapCache map[string]interface{}
}

func (this ExpressionEngineLexerMapCache) New() ExpressionEngineLexerMapCache {
	if this.mapCache == nil {
		this.mapCache = make(map[string]interface{})
	}
	return this
}

func (this *ExpressionEngineLexerMapCache) Set(expression string, lexer interface{}) error {
	if expression == "" {
		return utils.NewError("ExpressionEngineLexerMapCache", "set lexerMap chache key can not be ''!")
	}
	this.mapCache[expression] = lexer
	return nil
}
func (this *ExpressionEngineLexerMapCache) Get(expression string) (interface{}, error) {
	return this.mapCache[expression], nil
}

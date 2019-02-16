package GoMybatis

import (
	"github.com/zhuxiujia/GoMybatis/utils"
	"sync"
)

type ExpressionEngineLexerMapCache struct {
	mapCache map[string]interface{}
	lock     sync.RWMutex
}

func (it ExpressionEngineLexerMapCache) New() ExpressionEngineLexerMapCache {
	if it.mapCache == nil {
		it.mapCache = make(map[string]interface{})
	}
	return it
}

func (it *ExpressionEngineLexerMapCache) Set(expression string, lexer interface{}) error {
	if expression == "" {
		return utils.NewError("ExpressionEngineLexerMapCache", "set lexerMap chache key can not be ''!")
	}
	it.lock.Lock()
	defer it.lock.Unlock()
	it.mapCache[expression] = lexer
	return nil
}
func (it *ExpressionEngineLexerMapCache) Get(expression string) (interface{}, error) {
	var result interface{}
	it.lock.RLock()
	defer it.lock.RUnlock()
	result = it.mapCache[expression]
	return result, nil
}
func (it *ExpressionEngineLexerMapCache) Name() string  {
	return "ExpressionEngineLexerMapCache"
}

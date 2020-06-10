package GoMybatis

import (
	"github.com/zhuxiujia/GoMybatis/utils"
	"sync"
)

type ExpressionEngineLexerMapCache struct {
	mapCache sync.Map
	lock     sync.RWMutex
}

func (it ExpressionEngineLexerMapCache) New() ExpressionEngineLexerMapCache {
	return it
}

func (it *ExpressionEngineLexerMapCache) Set(expression string, lexer interface{}) error {
	if expression == "" {
		return utils.NewError("ExpressionEngineLexerMapCache", "set lexerMap chache key can not be ''!")
	}
	it.lock.Lock()
	defer it.lock.Unlock()
	it.mapCache.Store(expression, lexer)
	return nil
}
func (it *ExpressionEngineLexerMapCache) Get(expression string) (interface{}, error) {
	var result interface{}
	it.lock.RLock()
	defer it.lock.RUnlock()
	result, ok := it.mapCache.Load(expression)
	if !ok {
		return nil, nil
	}
	return result, nil
}
func (it *ExpressionEngineLexerMapCache) Name() string {
	return "ExpressionEngineLexerMapCache"
}

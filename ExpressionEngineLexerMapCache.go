package GoMybatis

import (
	"github.com/zhuxiujia/GoMybatis/utils"
	"sync"
)

type ExpressionEngineLexerMapCache struct {
	mapCache map[string]interface{}
	lock     sync.RWMutex
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
	this.lock.Lock()
	this.mapCache[expression] = lexer
	this.lock.Unlock()
	return nil
}
func (this *ExpressionEngineLexerMapCache) Get(expression string) (interface{}, error) {
	var result interface{}
	this.lock.RLock()
	result = this.mapCache[expression]
	this.lock.RUnlock()
	return result, nil
}

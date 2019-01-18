package GoMybatis

import (
	"encoding/json"
	"github.com/zhuxiujia/GoMybatis/lib/github.com/nytlabs/gojee"
	"github.com/zhuxiujia/GoMybatis/utils"
)

type ExpressionOperation = int

const (
	Operation_JSON_Marshal_UnMarshal_Map ExpressionOperation = iota //序列化和反序列化为json
	Operation_JSON_Unmarshal_Byte
)

//ExpressionEngineJee 是一个基于json表达式操作的第三方库实现
type ExpressionEngineJee struct {
}

func (this *ExpressionEngineJee) Lexer(lexerArg string) (interface{}, error) {
	tokenized, err := jee.Lexer(lexerArg)
	if err != nil {
		return nil, utils.MakeErrors("[GoMybatis][JeeExpressionEngine]", err.Error())
	}
	tree, err := jee.Parser(tokenized)
	if err != nil {
		return nil, utils.MakeErrors("[GoMybatis][JeeExpressionEngine]", err.Error())
	}
	return tree, nil
}

func (this *ExpressionEngineJee) Eval(compileResult interface{}, arg interface{}, operation int) (interface{}, error) {
	var jeeMsg jee.BMsg
	switch operation {
	case Operation_JSON_Marshal_UnMarshal_Map:
		//to json
		bytes, err := json.Marshal(arg.(map[string]interface{}))
		if err != nil {
			return nil, utils.MakeErrors("[GoMybatis][JeeExpressionEngine]", err.Error())
		}
		err = json.Unmarshal(bytes, &jeeMsg)
		if err != nil {
			return nil, utils.MakeErrors("[GoMybatis][JeeExpressionEngine]", err.Error())
		}
		break
	case Operation_JSON_Unmarshal_Byte:
		//to json
		err := json.Unmarshal(arg.([]byte), &jeeMsg)
		if err != nil {
			return nil, utils.MakeErrors("[GoMybatis][JeeExpressionEngine]", err.Error())
		}
		break

	}
	result, err := jee.Eval(compileResult.(*jee.TokenTree), jeeMsg)
	if err != nil {
		return nil, utils.MakeErrors("[GoMybatis][JeeExpressionEngine]", err.Error())
	}
	return result, nil
}

package GoMybatis

import (
	"bytes"
	"encoding/json"
	"github.com/zhuxiujia/GoMybatis/lib/github.com/nytlabs/gojee"
	"github.com/zhuxiujia/GoMybatis/utils"
	"strings"
)

type ExpressionOperation = int

const (
	JeeOperation_Marshal_Map ExpressionOperation = iota //序列化和反序列化为json
	JeeOperation_Unmarshal_Byte
)

//ExpressionEngineJee 是一个基于json表达式操作的第三方库实现
type ExpressionEngineJee struct {
}

//编译一个表达式
//参数：lexerArg 表达式内容
//返回：interface{} 编译结果,error 错误
func (this *ExpressionEngineJee) Lexer(lexerArg string) (interface{}, error) {
	tokenized, err := jee.Lexer(this.LexerAndOrSupport(lexerArg))
	if err != nil {
		return nil, utils.NewError("ExpressionEngineJee", err)
	}
	tree, err := jee.Parser(tokenized)
	if err != nil {
		return nil, utils.NewError("ExpressionEngineJee", err)
	}
	return tree, nil
}
//执行一个表达式
//参数：lexerResult=编译结果，arg=参数
//返回：执行结果，错误
func (this *ExpressionEngineJee) Eval(compileResult interface{}, arg interface{}, operation int) (interface{}, error) {
	var jeeMsg jee.BMsg
	switch operation {
	case JeeOperation_Marshal_Map:
		//to json，针对arg是map[string]interface{}的数据类型
		bytes, err := json.Marshal(arg.(map[string]interface{}))
		if err != nil {
			return nil, utils.NewError("ExpressionEngineJee", err)
		}
		err = json.Unmarshal(bytes, &jeeMsg)
		if err != nil {
			return nil, utils.NewError("ExpressionEngineJee", err)
		}
		break
	case JeeOperation_Unmarshal_Byte:
		//to json,针对arg是json byte的数据类型
		err := json.Unmarshal(arg.([]byte), &jeeMsg)
		if err != nil {
			return nil, utils.NewError("ExpressionEngineJee", err)
		}
		break
	default:
		return nil, utils.NewError("ExpressionEngineJee", "unknow operation value = ",operation,"!")
	}
	result, err := jee.Eval(compileResult.(*jee.TokenTree), jeeMsg)
	if err != nil {
		return nil, utils.NewError("ExpressionEngineJee", err)
	}
	return result, nil
}

//编译后立即执行
func (this *ExpressionEngineJee) LexerEval(lexerArg string, arg interface{}, operation int) (interface{}, error) {
	var lexer, error = this.Lexer(lexerArg)
	if error != nil {
		return nil, error
	}
	return this.Eval(lexer, arg, operation)
}

//添加and 和 or 语法支持
func (this *ExpressionEngineJee) LexerAndOrSupport(lexerArg string) string {
	var buf bytes.Buffer
	strs := strings.Split(lexerArg, " or ")
	if len(strs) > 1 {
		buf.Reset()
		var buf bytes.Buffer
		var len = len(strs)
		for index, k := range strs {
			buf.WriteString("(")
			buf.WriteString(k)
			buf.WriteString(")")
			if index+1 < len {
				buf.WriteString(" || ")
			}
		}
		lexerArg = buf.String()
		buf.Reset()
	}
	strs = strings.Split(lexerArg, " and ")
	if len(strs) > 1 {
		var len = len(strs)
		for index, k := range strs {
			buf.WriteString("(")
			buf.WriteString(k)
			buf.WriteString(")")
			if index+1 < len {
				buf.WriteString(" && ")
			}
		}
		lexerArg = buf.String()
		buf.Reset()
	}
	return lexerArg
}

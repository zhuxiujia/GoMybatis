package GoMybatis

import (
	"bytes"
	"github.com/zhuxiujia/GoMybatis/utils"
	"strings"
)

type SqlNodeType int

const (
	NArg    SqlNodeType = iota
	NString             //string 节点
	NNil                //空节点
	NBinary             //二元计算节点
	NOpt                //操作符节点

	NIf
	NTrim
)

func (it SqlNodeType) ToString() string {
	switch it {
	case NString:
		return "NString"
	case NNil:
		return "NNil"
	case NBinary:
		return "NBinary"
	case NOpt:
		return "NOpt"
	}
	return "Unknow"
}

type SqlNode interface {
	Type() SqlNodeType
	Eval(env map[string]interface{}) (interface{}, error)
}

type StringNode struct {
	value string
	t     SqlNodeType
}

func (it StringNode) Type() SqlNodeType {
	return NString
}

func (it StringNode) Eval(env map[string]interface{}) (interface{}, error) {
	var sqlArgTypeConvert = env["SqlArgTypeConvert"]
	var expressionEngineProxy = env["*ExpressionEngineProxy"]

	var convert SqlArgTypeConvert
	var proxy *ExpressionEngineProxy
	if sqlArgTypeConvert != nil {
		convert = sqlArgTypeConvert.(SqlArgTypeConvert)
	}
	if expressionEngineProxy != nil {
		proxy = expressionEngineProxy.(*ExpressionEngineProxy)
	}
	var result, e = replaceArg(it.value, env, convert, proxy)
	return result, e
}

type IfNode struct {
	childs []SqlNode
	test   string
	t      SqlNodeType
}

func (it IfNode) Type() SqlNodeType {
	return NIf
}

func (it IfNode) Eval(env map[string]interface{}) (interface{}, error) {
	var expressionEngineProxy = env["*ExpressionEngineProxy"]
	var proxy *ExpressionEngineProxy
	if expressionEngineProxy != nil {
		proxy = expressionEngineProxy.(*ExpressionEngineProxy)
	}
	var result, err = proxy.LexerAndEval(it.test, env)
	if err != nil {
		err = utils.NewError("GoMybatisSqlBuilder", "[GoMybatis] <test `", it.test, `> fail,`, err.Error())
	}
	if result.(bool) {
		return DoChildNodes(it.childs, env)
	}
	return nil, nil
}

type TrimNode struct {
	childs          []SqlNode
	prefix          string
	suffix          string
	suffixOverrides string
	prefixOverrides string
	t               SqlNodeType
}

func (it TrimNode) Type() SqlNodeType {
	return NTrim
}

func (it TrimNode) Eval(env map[string]interface{}) (interface{}, error) {
	var sql, err = DoChildNodes(it.childs, env)
	if err != nil {
		return nil, err
	}
	var tempTrimSqlString = strings.Trim(sql.(string), " ")
	if it.prefixOverrides != "" {
		var prefixOverridesArray = strings.Split(it.prefixOverrides, "|")
		if len(prefixOverridesArray) > 0 {
			for _, v := range prefixOverridesArray {
				tempTrimSqlString = strings.TrimPrefix(tempTrimSqlString, v)
			}
		}
	}
	if it.suffixOverrides != "" {
		var suffixOverrideArray = strings.Split(it.suffixOverrides, "|")
		if len(suffixOverrideArray) > 0 {
			for _, v := range suffixOverrideArray {
				tempTrimSqlString = strings.TrimSuffix(tempTrimSqlString, v)
			}
		}
	}
	var newBuffer bytes.Buffer
	newBuffer.WriteString(` `)
	newBuffer.WriteString(it.prefix)
	newBuffer.WriteString(` `)
	newBuffer.WriteString(tempTrimSqlString)
	newBuffer.WriteString(` `)
	newBuffer.WriteString(it.suffix)
	return newBuffer.String(), nil
}

//执行子所有节点
func DoChildNodes(childs []SqlNode, env map[string]interface{}) (interface{}, error) {
	if childs == nil {
		return nil, nil
	}
	var sql bytes.Buffer
	for _, v := range childs {
		var r, e = v.Eval(env)
		if e != nil {
			return nil, e
		}
		if r != nil {
			sql.WriteString(r.(string))
		}
	}
	return sql.String(), nil
}

////计算节点
//type ArrayNode struct {
//	childs []SqlNode
//	t      SqlNodeType
//}
//
//func (it ArrayNode) Type() SqlNodeType {
//	return NArray
//}
//
//func (it ArrayNode) Eval(env map[string]interface{}) (interface{}, error) {
//	var sql bytes.Buffer
//	if it.childs != nil {
//		for _, v := range it.childs {
//			var r, e = v.Eval(env)
//			if e != nil {
//				return nil, e
//			}
//			if r != nil {
//				sql.WriteString(r.(string))
//			}
//		}
//	}
//	return sql.String(), nil
//}

package GoMybatis

import (
	"bytes"
	"github.com/zhuxiujia/GoMybatis/utils"
	"reflect"
	"strings"
)

//sql构建抽象语法树节点
type SqlNode interface {
	Type() SqlNodeType
	Eval(env map[string]interface{}) (*bytes.Buffer, error)
}

//字符串节点
type StringNode struct {
	value string
	t     SqlNodeType
}

func (it *StringNode) Type() SqlNodeType {
	return NString
}

func (it *StringNode) Eval(env map[string]interface{}) (*bytes.Buffer, error) {
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
	var r, e = replaceArg(it.value, env, convert, proxy)
	if e != nil {
		return nil, e
	}
	var buf bytes.Buffer
	buf.WriteString(r)
	return &buf, nil
}

//判断节点
type IfNode struct {
	childs []SqlNode
	test   string
	t      SqlNodeType
}

func (it *IfNode) Type() SqlNodeType {
	return NIf
}

func (it *IfNode) Eval(env map[string]interface{}) (*bytes.Buffer, error) {
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

//Trim操作节点
type TrimNode struct {
	childs          []SqlNode
	prefix          string
	suffix          string
	suffixOverrides string
	prefixOverrides string
	t               SqlNodeType
}

func (it *TrimNode) Type() SqlNodeType {
	return NTrim
}

func (it *TrimNode) Eval(env map[string]interface{}) (*bytes.Buffer, error) {
	var sql, err = DoChildNodes(it.childs, env)
	if err != nil {
		return nil, err
	}
	if sql == nil {
		return nil, nil
	}
	var tempTrimSqlString = strings.Trim(sql.String(), " ")
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
	return &newBuffer, nil
}

//set节点
type SetNode struct {
	childs []SqlNode
	t      SqlNodeType
}

func (it *SetNode) Type() SqlNodeType {
	return NSet
}

func (it *SetNode) Eval(env map[string]interface{}) (*bytes.Buffer, error) {
	var sql, err = DoChildNodes(it.childs, env)
	if err != nil {
		return nil, err
	}
	if sql == nil {
		return nil, nil
	}
	var trim bytes.Buffer
	if sql != nil {
		var trimString = strings.Trim(sql.String(), DefaultOverrides)
		trim.Reset()
		trim.WriteString(` `)
		trim.WriteString(` set `)
		trim.WriteString(trimString)
		trim.WriteString(` `)
	}
	return &trim, nil
}

//foreach 节点
type ForEachNode struct {
	childs []SqlNode
	t      SqlNodeType

	collection string
	index      string
	item       string
	open       string
	close      string
	separator  string
}

func (it *ForEachNode) Type() SqlNodeType {
	return NForEach
}

func (it *ForEachNode) Eval(env map[string]interface{}) (*bytes.Buffer, error) {
	if it.collection == "" {
		panic(`[GoMybatis] collection value can not be "" in <foreach collection=""> !`)
	}
	var tempSql bytes.Buffer
	var datas = env[it.collection]
	var collectionValue = reflect.ValueOf(datas)
	if collectionValue.Kind() != reflect.Slice && collectionValue.Kind() != reflect.Map {
		panic(`[GoMybatis] collection value must be a slice or map !`)
	}
	var collectionValueLen = collectionValue.Len()
	if collectionValueLen == 0 {
		return nil, nil
	}
	if it.index == "" {
		it.index = "index"
	}
	if it.item == "" {
		it.item = "item"
	}
	switch collectionValue.Kind() {
	case reflect.Map:
		var mapKeys = collectionValue.MapKeys()
		var collectionKeyLen = len(mapKeys)
		if collectionKeyLen == 0 {
			return nil, nil
		}
		var tempArgMap = env
		for _, keyValue := range mapKeys {
			var key = keyValue.Interface()
			var collectionItem = collectionValue.MapIndex(keyValue)
			if it.item != "" {
				tempArgMap[it.item] = collectionItem.Interface()
			}
			tempArgMap[it.index] = key
			var r, err = DoChildNodes(it.childs, tempArgMap)
			if err != nil {
				return nil, err
			}
			if r != nil {
				tempSql.WriteString(r.String())
			}
			tempSql.WriteString(it.separator)
			delete(tempArgMap, it.item)
		}
		break
	case reflect.Slice:
		var tempArgMap = env
		for i := 0; i < collectionValueLen; i++ {
			var collectionItem = collectionValue.Index(i)
			if it.item != "" {
				tempArgMap[it.item] = collectionItem.Interface()
			}
			if it.index != "" {
				tempArgMap[it.index] = i
			}
			var r, err = DoChildNodes(it.childs, tempArgMap)
			if err != nil {
				return nil, err
			}
			if r != nil {
				tempSql.WriteString(r.String())
			}
			tempSql.WriteString(it.separator)
			delete(tempArgMap, it.item)
		}
		break
	}
	var newTempSql bytes.Buffer
	var tempSqlString = strings.Trim(strings.Trim(tempSql.String(), " "), it.separator)
	newTempSql.WriteString(it.open)
	newTempSql.WriteString(tempSqlString)
	newTempSql.WriteString(it.close)
	tempSql.Reset()
	return &newTempSql, nil
}

//执行子所有节点
func DoChildNodes(childs []SqlNode, env map[string]interface{}) (*bytes.Buffer, error) {
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
			sql.WriteString(r.String())
		}
	}
	return &sql, nil
}

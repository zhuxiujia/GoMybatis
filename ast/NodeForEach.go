package ast

import (
	"bytes"
	"errors"
	"github.com/zhuxiujia/GoMybatis/stmt"
	"reflect"
)

//foreach 节点
type NodeForEach struct {
	childs []Node
	t      NodeType

	collection string
	index      string
	item       string
	open       string
	close      string
	separator  string

	holder *NodeConfigHolder
}

func (it *NodeForEach) Type() NodeType {
	return NForEach
}

func (it *NodeForEach) Eval(env map[string]interface{}, arg_array *[]interface{}, stmtConvert stmt.StmtIndexConvert) ([]byte, error) {
	if it.collection == "" {
		panic(`[GoMybatis] collection value can not be "" in <foreach collection=""> !`)
	}
	var tempSql bytes.Buffer

	//exec lexer
	var err error
	evalData, err := it.holder.GetExpressionEngineProxy().LexerAndEval(it.collection, env)
	if err != nil {
		return nil, errors.New(it.holder.GetExpressionEngineProxy().Name() + ":" + err.Error())
	}
	var collectionValue = reflect.ValueOf(evalData)
	var kind = collectionValue.Kind()
	if kind == reflect.Invalid {
		return nil, errors.New(it.holder.GetExpressionEngineProxy().Name() + ": collection value is invalid value!")
	}
	if kind != reflect.Slice && kind != reflect.Array && kind != reflect.Map {
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
			var r, err = DoChildNodes(it.childs, tempArgMap, arg_array, stmtConvert)
			if err != nil {
				return nil, err
			}
			if r != nil {
				tempSql.Write(r)
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
			var r, err = DoChildNodes(it.childs, tempArgMap, arg_array, stmtConvert)
			if err != nil {
				return nil, err
			}
			if r != nil {
				tempSql.Write(r)
			}
			tempSql.WriteString(it.separator)
			delete(tempArgMap, it.item)
		}
		break
	}
	var newTempSql bytes.Buffer
	var tempSqlString = bytes.Trim(tempSql.Bytes(), it.separator)
	tempSql.Reset()
	newTempSql.WriteString(it.open)
	newTempSql.Write(tempSqlString)
	newTempSql.WriteString(it.close)
	var newTempSqlBytes = newTempSql.Bytes()
	newTempSql.Reset()
	return newTempSqlBytes, nil
}

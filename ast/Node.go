package ast

import (
	"bytes"
	"github.com/zhuxiujia/GoMybatis/stmt"
)

//sql构建抽象语法树节点
type Node interface {
	Type() NodeType
	Eval(env map[string]interface{}, arg_array *[]interface{}, stmtConvert stmt.StmtIndexConvert) ([]byte, error)
}

//执行子所有节点
func DoChildNodes(childNodes []Node, env map[string]interface{}, arg_array *[]interface{}, stmtConvert stmt.StmtIndexConvert) ([]byte, error) {
	if childNodes == nil {
		return nil, nil
	}
	var sql bytes.Buffer
	for _, v := range childNodes {
		var r, e = v.Eval(env, arg_array, stmtConvert)
		if e != nil {
			return nil, e
		}
		if r != nil {
			sql.Write(r)
		}
	}
	var bytes = sql.Bytes()
	sql.Reset()
	return bytes, nil
}

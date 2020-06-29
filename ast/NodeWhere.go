package ast

import (
	"bytes"
	"github.com/zhuxiujia/GoMybatis/stmt"
)

//Trim操作节点
type NodeWhere struct {
	childs []Node
	t      NodeType
}

func (it *NodeWhere) Type() NodeType {
	return NWhere
}

func (it *NodeWhere) Eval(env map[string]interface{}, arg_array *[]interface{}, stmtConvert stmt.StmtIndexConvert) ([]byte, error) {
	var sql, err = DoChildNodes(it.childs, env, arg_array, stmtConvert)
	if err != nil {
		return nil, err
	}
	if sql == nil {
		return nil, nil
	}
	for {
		if bytes.HasPrefix(sql, []byte(" ")) {
			sql = bytes.Trim(sql, " ")
		} else {
			break
		}
	}
	if len(sql) == 0 {
		return sql, nil
	}

	sql = bytes.TrimPrefix(sql, []byte("and"))
	sql = bytes.TrimPrefix(sql, []byte("AND"))
	sql = bytes.TrimPrefix(sql, []byte("And"))

	sql = bytes.TrimPrefix(sql, []byte("or"))
	sql = bytes.TrimPrefix(sql, []byte("OR"))
	sql = bytes.TrimPrefix(sql, []byte("Or"))

	var newBuffer bytes.Buffer
	newBuffer.WriteString(` `)
	newBuffer.WriteString("WHERE")
	newBuffer.WriteString(` `)
	newBuffer.Write(sql)
	newBuffer.WriteString(` `)

	var newBufferBytes = newBuffer.Bytes()
	newBuffer.Reset()
	return newBufferBytes, nil
}

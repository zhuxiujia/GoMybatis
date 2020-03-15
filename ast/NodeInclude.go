package ast

import "github.com/zhuxiujia/GoMybatis/stmt"

type NodeInclude struct {
	childs []Node
	t      NodeType
}

func (it *NodeInclude) Type() NodeType {
	return NInclude
}

func (it *NodeInclude) Eval(env map[string]interface{}, arg_array *[]interface{}, stmtConvert stmt.StmtIndexConvert) ([]byte, error) {
	var sql, err = DoChildNodes(it.childs, env, arg_array, stmtConvert)
	return sql, err
}

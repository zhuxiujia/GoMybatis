package ast

import "github.com/zhuxiujia/GoMybatis/stmt"

type NodeOtherwise struct {
	childs []Node
	t      NodeType
}

func (it *NodeOtherwise) Type() NodeType {
	return NOtherwise
}

func (it *NodeOtherwise) Eval(env map[string]interface{}, arg_array *[]interface{}, stmtConvert stmt.StmtIndexConvert) ([]byte, error) {
	var r, e = DoChildNodes(it.childs, env, arg_array, stmtConvert)
	if e != nil {
		return nil, e
	}
	return r, nil
}

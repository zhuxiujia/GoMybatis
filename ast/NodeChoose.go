package ast

import "github.com/zhuxiujia/GoMybatis/stmt"

type NodeChoose struct {
	t             NodeType
	whenNodes     []Node
	otherwiseNode Node
}

func (it *NodeChoose) Type() NodeType {
	return NChoose
}

func (it *NodeChoose) Eval(env map[string]interface{}, arg_array *[]interface{}, stmtConvert stmt.StmtIndexConvert) ([]byte, error) {
	if it.whenNodes == nil && it.otherwiseNode == nil {
		return nil, nil
	}
	for _, v := range it.whenNodes {
		var r, e = v.Eval(env, arg_array, stmtConvert)
		if e != nil {
			return nil, e
		}
		if r != nil {
			return r, nil
		}
	}
	return it.otherwiseNode.Eval(env, arg_array, stmtConvert)
}

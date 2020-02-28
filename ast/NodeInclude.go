package ast

type NodeInclude struct {
	childs []Node
	t      NodeType
}

func (it *NodeInclude) Type() NodeType {
	return NInclude
}

func (it *NodeInclude) Eval(env map[string]interface{}, arg_array *[]interface{}) ([]byte, error) {
	var sql, err = DoChildNodes(it.childs, env, arg_array)
	return sql, err
}

package GoMybatis

type NodeInclude struct {
	childs []Node
	t      NodeType
}

func (it *NodeInclude) Type() NodeType {
	return NInclude
}

func (it *NodeInclude) Eval(env map[string]interface{}) ([]byte, error) {
	var sql, err = DoChildNodes(it.childs, env)
	return sql, err
}

package GoMybatis


type NodeOtherwise struct {
	childs []Node
	t      NodeType
}

func (it *NodeOtherwise) Type() NodeType {
	return NOtherwise
}

func (it *NodeOtherwise) Eval(env map[string]interface{}) ([]byte, error) {
	var r, e = DoChildNodes(it.childs, env)
	if e != nil {
		return nil, e
	}
	return r, nil
}


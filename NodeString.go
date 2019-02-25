package GoMybatis


//字符串节点
type NodeString struct {
	value string
	t     NodeType
}

func (it *NodeString) Type() NodeType {
	return NString
}

func (it *NodeString) Eval(env map[string]interface{}) ([]byte, error) {
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
	return []byte(r), nil
}
package GoMybatis

//字符串节点
type NodeString struct {
	value string
	t     NodeType

	//args
	expressMap          map[string]int //express表 key：name
	noConvertExpressMap map[string]int
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

	var data = it.value
	var err error
	if it.expressMap != nil {
		data, err = Replace(`#{`, it.expressMap, data, convert, env, proxy)
		if err != nil {
			return nil, err
		}
	}
	if it.noConvertExpressMap != nil {
		data, err = Replace(`${`, it.noConvertExpressMap, data, convert, env, proxy)
		if err != nil {
			return nil, err
		}
	}
	return []byte(data), nil
}

package GoMybatis

type NodeConfigHolder struct {
	convert SqlArgTypeConvert
	proxy   *ExpressionEngineProxy
}

func (it *NodeConfigHolder) GetSqlArgTypeConvert() SqlArgTypeConvert {
	return it.convert
}

func (it *NodeConfigHolder) GetExpressionEngineProxy() *ExpressionEngineProxy {
	return it.proxy
}

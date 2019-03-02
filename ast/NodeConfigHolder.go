package ast


type NodeConfigHolder struct {
	Convert SqlArgTypeConvert
	Proxy   ExpressionEngine
}

func (it *NodeConfigHolder) GetSqlArgTypeConvert() SqlArgTypeConvert {
	return it.Convert
}

func (it *NodeConfigHolder) GetExpressionEngineProxy() ExpressionEngine {
	return it.Proxy
}

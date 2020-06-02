package ast

type NodeConfigHolder struct {
	Proxy ExpressionEngine
}

func (it *NodeConfigHolder) GetExpressionEngineProxy() ExpressionEngine {
	return it.Proxy
}

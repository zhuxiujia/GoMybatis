package GoMybatis

import "github.com/zhuxiujia/GoMybatis/utils"

type NodeWhen struct {
	childs []Node
	test   string
	t      NodeType

	expressionEngineProxy *ExpressionEngineProxy
}

func (it *NodeWhen) Type() NodeType {
	return NWhen
}

func (it *NodeWhen) Eval(env map[string]interface{}) ([]byte, error) {
	var expressionEngineProxy = env["*ExpressionEngineProxy"]
	var proxy *ExpressionEngineProxy
	if expressionEngineProxy != nil {
		proxy = expressionEngineProxy.(*ExpressionEngineProxy)
	}
	var result, err = proxy.LexerAndEval(it.test, env)
	if err != nil {
		err = utils.NewError("GoMybatisSqlBuilder", "[GoMybatis] <test `", it.test, `> fail,`, err.Error())
	}
	if result.(bool) {
		return DoChildNodes(it.childs, env)
	}
	return nil, nil
}


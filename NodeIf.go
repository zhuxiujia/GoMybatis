package GoMybatis

import "github.com/zhuxiujia/GoMybatis/utils"

//判断节点
type NodeIf struct {
	childs []Node
	t      NodeType
	test   string
}

func (it *NodeIf) Type() NodeType {
	return NIf
}

func (it *NodeIf) Eval(env map[string]interface{}) ([]byte, error) {
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

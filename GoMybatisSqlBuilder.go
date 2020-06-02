package GoMybatis

import (
	"github.com/zhuxiujia/GoMybatis/ast"
	"github.com/zhuxiujia/GoMybatis/stmt"
)

type GoMybatisSqlBuilder struct {
	expressionEngineProxy ExpressionEngineProxy
	enableLog             bool
	nodeParser            ast.NodeParser
}

func (it *GoMybatisSqlBuilder) ExpressionEngineProxy() *ExpressionEngineProxy {
	return &it.expressionEngineProxy
}

func (it GoMybatisSqlBuilder) New(expressionEngine ExpressionEngineProxy, log Log, enableLog bool) GoMybatisSqlBuilder {
	it.expressionEngineProxy = expressionEngine
	it.enableLog = enableLog
	it.nodeParser = ast.NodeParser{
		Holder: ast.NodeConfigHolder{
			Proxy: &expressionEngine,
		},
	}
	return it
}

func (it *GoMybatisSqlBuilder) BuildSql(paramMap map[string]interface{}, nodes []ast.Node, arg_array *[]interface{}, stmtConvert stmt.StmtIndexConvert) (string, error) {
	//抽象语法树节点构建
	var sql, err = ast.DoChildNodes(nodes, paramMap, arg_array, stmtConvert)
	if err != nil {
		return "", err
	}
	var sqlStr = string(sql)
	return sqlStr, nil
}

func (it *GoMybatisSqlBuilder) SetEnableLog(enable bool) {
	it.enableLog = enable
}
func (it *GoMybatisSqlBuilder) EnableLog() bool {
	return it.enableLog
}

func (it *GoMybatisSqlBuilder) NodeParser() ast.NodeParser {
	return it.nodeParser
}

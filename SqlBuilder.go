package GoMybatis

import "github.com/zhuxiujia/GoMybatis/ast"

//sql文本构建
type SqlBuilder interface {
	BuildSql(paramMap map[string]interface{}, nodes []ast.Node) (string, error)
	ExpressionEngineProxy() *ExpressionEngineProxy
	SqlArgTypeConvert() ast.SqlArgTypeConvert
	LogSystem() *LogSystem
	SetEnableLog(enable bool)
	EnableLog() bool
	NodeParser() ast.NodeParser
}

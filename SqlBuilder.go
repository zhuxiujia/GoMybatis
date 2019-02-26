package GoMybatis

//sql文本构建
type SqlBuilder interface {
	BuildSql(paramMap map[string]interface{}, nodes []Node) (string, error)
	ExpressionEngineProxy() *ExpressionEngineProxy
	SqlArgTypeConvert() SqlArgTypeConvert
	LogSystem() *LogSystem
	SetEnableLog(enable bool)
	EnableLog() bool
	NodeParser() NodeParser
}

package GoMybatis

//sql文本构建
type SqlBuilder interface {
	BuildSql(paramMap map[string]interface{}, mapperXml *MapperXml) (string, error)
	ExpressionEngineProxy() *ExpressionEngineProxy
	SqlArgTypeConvert() SqlArgTypeConvert
	LogSystem() *LogSystem
	SetEnableLog(enable bool)
	EnableLog() bool
}

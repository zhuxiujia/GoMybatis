package GoMybatis

//sql文本构建
type SqlBuilder interface {
	BuildSql(paramMap map[string]SqlArg, mapperXml *MapperXml) (string, error)
	ExpressionEngineProxy() ExpressionEngineProxy
	SqlArgTypeConvert() SqlArgTypeConvert
	ExpressionTypeConvert() ExpressionTypeConvert
	LogSystem() *LogSystem
}

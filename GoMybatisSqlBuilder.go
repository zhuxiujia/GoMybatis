package GoMybatis

import (
	"time"
)

type GoMybatisSqlBuilder struct {
	sqlArgTypeConvert     SqlArgTypeConvert
	expressionEngineProxy ExpressionEngineProxy
	logSystem             *LogSystem
	enableLog             bool

	nodeParser NodeParser
}

func (it *GoMybatisSqlBuilder) ExpressionEngineProxy() *ExpressionEngineProxy {
	return &it.expressionEngineProxy
}
func (it *GoMybatisSqlBuilder) SqlArgTypeConvert() SqlArgTypeConvert {
	return it.sqlArgTypeConvert
}

func (it GoMybatisSqlBuilder) New(SqlArgTypeConvert SqlArgTypeConvert, expressionEngine ExpressionEngineProxy, log Log, enableLog bool) GoMybatisSqlBuilder {
	it.sqlArgTypeConvert = SqlArgTypeConvert
	it.expressionEngineProxy = expressionEngine
	it.enableLog = enableLog
	if enableLog {
		var logSystem, err = LogSystem{}.New(log, log.QueueLen())
		if err != nil {
			panic(err)
		}
		it.logSystem = &logSystem
	}
	it.nodeParser = NodeParser{
		holder: NodeConfigHolder{
			convert: SqlArgTypeConvert,
			proxy:   &expressionEngine,
		},
	}
	return it
}

func (it *GoMybatisSqlBuilder) BuildSql(paramMap map[string]interface{}, nodes []Node) (string, error) {
	//抽象语法树节点构建
	var sql, err = DoChildNodes(nodes, paramMap)
	if err != nil {
		return "", err
	}
	var sqlStr = string(sql)
	if it.enableLog {
		var now, _ = time.Now().MarshalText()
		it.logSystem.SendLog("[GoMybatis] [", string(now), "] Preparing sql ==> ", sqlStr)
	}
	return sqlStr, nil
}

func (it *GoMybatisSqlBuilder) SetEnableLog(enable bool) {
	it.enableLog = enable
}
func (it *GoMybatisSqlBuilder) EnableLog() bool {
	return it.enableLog
}

func (it *GoMybatisSqlBuilder) LogSystem() *LogSystem {
	return it.logSystem
}

func (it *GoMybatisSqlBuilder) NodeParser() NodeParser {
	return it.nodeParser
}

package GoMybatis

import (
	"github.com/zhuxiujia/GoMybatis/ast"
	"time"
)

type GoMybatisSqlBuilder struct {
	sqlArgTypeConvert     ast.SqlArgTypeConvert
	expressionEngineProxy ExpressionEngineProxy
	logSystem             *LogSystem
	enableLog             bool

	nodeParser ast.NodeParser
}

func (it *GoMybatisSqlBuilder) ExpressionEngineProxy() *ExpressionEngineProxy {
	return &it.expressionEngineProxy
}
func (it *GoMybatisSqlBuilder) SqlArgTypeConvert() ast.SqlArgTypeConvert {
	return it.sqlArgTypeConvert
}

func (it GoMybatisSqlBuilder) New(SqlArgTypeConvert ast.SqlArgTypeConvert, expressionEngine ExpressionEngineProxy, log Log, enableLog bool) GoMybatisSqlBuilder {
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
	it.nodeParser = ast.NodeParser{
		Holder: ast.NodeConfigHolder{
			Convert: SqlArgTypeConvert,
			Proxy:   &expressionEngine,
		},
	}
	return it
}

func (it *GoMybatisSqlBuilder) BuildSql(paramMap map[string]interface{}, nodes []ast.Node) (string, error) {
	//抽象语法树节点构建
	var sql, err = ast.DoChildNodes(nodes, paramMap)
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

func (it *GoMybatisSqlBuilder) NodeParser() ast.NodeParser {
	return it.nodeParser
}

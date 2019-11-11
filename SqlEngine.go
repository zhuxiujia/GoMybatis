package GoMybatis

import (
	"database/sql"
	"time"

	"github.com/zhuxiujia/GoMybatis/ast"
	"github.com/zhuxiujia/GoMybatis/tx"
)


type Result struct {
	LastInsertId int64
	RowsAffected int64
}

type SqlType string

func (s SqlType) newPoint() interface{} {
	switch s {
	//数值类型
	case "TINYINT",       //8(bit)
		"SMALLINT",       //16(bit)
		"MEDIUMINT",      //34(bit)
		"INT", "INTEGER": //32(bit)
		var i int
		return &i
	case "BIGINT": //64(bit)
		var i int
		return &i
	case "FLOAT": //32(bit)
		var i float32
		return &i
	case "DOUBLE", "DECIMAL":
		var i float64
		return &i
	//字符串类型
	case "CHAR",      //0-255(byte)
		"VARCHAR",    //0-65535(byte)
		"TINYBLOB",   //0-255(byte)
		"TINYTEXT",   //0-255(byte)
		"BLOB",       //0-65535(byte)
		"TEXT",       //0-65535(byte)
		"MEDIUMBLOB", //0-16777215(byte)
		"MEDIUMTEXT", //0-16777215(byte)
		"LONGBLOB",   //0-4294967295(byte)
		"LONGTEXT":   //0-4294967295(byte)
		var i string
		return &i
	//日期和时间类型
	case "DATE",     //YYYY-MM-DD
		"TIME",      //HH:MM:SS
		"YEAR",      //YYYY
		"DATETIME",  //YYYY-MM-DD HH:MM:SS
		"TIMESTAMP": //YYYY-MM-DD HH:MM:SS
		var i time.Time
		return &i
	}
	return nil
}

type QueryResult struct {
	data      []map[string][]byte
	columnMap map[string]SqlType
}

func (q QueryResult) Index(index int) map[string][]byte {
	return q.data[index]
}

func (q QueryResult) Rows() int {
	return len(q.data)
}

func (q QueryResult) IsBlank() bool {
	return q.data == nil || len(q.data) == 0
}

func (q *QueryResult) append(cell map[string][]byte) {
	q.data = append(q.data, cell)
}

func (q QueryResult) SqlType(field string) SqlType {
	return q.columnMap[field]
}

type Session interface {
	Id() string
	Query(sqlorArgs string) (QueryResult, error)
	Exec(sqlorArgs string) (*Result, error)
	Rollback() error
	Commit() error
	Begin(p *tx.Propagation) error
	Close()
	LastPROPAGATION() *tx.Propagation
}

//产生session的引擎
type SessionEngine interface {
	//打开数据库
	Open(driverName, dataSourceName string) (*sql.DB, error)
	//写方法到mapper
	WriteMapperPtr(ptr interface{}, xml []byte)
	//引擎名称
	Name() string
	//创建session
	NewSession(mapperName string) (Session, error)
	//获取数据源路由
	DataSourceRouter() DataSourceRouter
	//设置数据源路由
	SetDataSourceRouter(router DataSourceRouter)

	//是否启用日志
	LogEnable() bool

	//是否启用日志
	SetLogEnable(enable bool)

	//获取日志实现类
	Log() Log

	//设置日志实现类
	SetLog(log Log)

	//session工厂
	SessionFactory() *SessionFactory

	//设置session工厂
	SetSessionFactory(factory *SessionFactory)

	//sql类型转换器
	SqlArgTypeConvert() ast.SqlArgTypeConvert

	//设置sql类型转换器
	SetSqlArgTypeConvert(convert ast.SqlArgTypeConvert)

	//表达式执行引擎
	ExpressionEngine() ast.ExpressionEngine

	//设置表达式执行引擎
	SetExpressionEngine(engine ast.ExpressionEngine)

	//sql构建器
	SqlBuilder() SqlBuilder

	//设置sql构建器
	SetSqlBuilder(builder SqlBuilder)

	//sql查询结果解析器
	SqlResultDecoder() SqlResultDecoder

	//设置sql查询结果解析器
	SetSqlResultDecoder(decoder SqlResultDecoder)

	//模板解析器
	TempleteDecoder() TempleteDecoder

	//设置模板解析器
	SetTempleteDecoder(decoder TempleteDecoder)

	RegisterObj(ptr interface{}, name string)

	GetObj(name string) interface{}

	//（注意（该方法需要在多协程环境下调用）启用会从栈获取协程id，有一定性能消耗，换取最大的事务定义便捷）
	GoroutineSessionMap() *GoroutineSessionMap

	//是否启用goroutineIDEnable（注意（该方法需要在多协程环境下调用）启用会从栈获取协程id，有一定性能消耗，换取最大的事务定义便捷）
	SetGoroutineIDEnable(enable bool)

	//是否启用goroutineIDEnable（注意（该方法需要在多协程环境下调用）启用会从栈获取协程id，有一定性能消耗，换取最大的事务定义便捷）
	GoroutineIDEnable() bool

	LogSystem() *LogSystem
}

package GoMybatis

import "database/sql"

type Result struct {
	LastInsertId int64
	RowsAffected int64
}

type Session interface {
	Id() string
	Query(sqlorArgs string) ([]map[string][]byte, error)
	Exec(sqlorArgs string) (*Result, error)
	Rollback() error
	Commit() error
	Begin() error
	Close()
}

//产生session的引擎
type SessionEngine interface {
	//写方法到mapper
	WriteMapperPtr(ptr interface{}, xml []byte)
	//引擎名称
	Name() string
	//创建session
	NewSession(mapperName string) (Session, error)
	//获取数据库map
	DBMap() map[string]*sql.DB
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

	//表达式数据类型转换器
	ExpressionTypeConvert() ExpressionTypeConvert

	//设置表达式数据类型转换器
	SetExpressionTypeConvert(convert ExpressionTypeConvert)

	//sql类型转换器
	SqlArgTypeConvert() SqlArgTypeConvert

	//设置sql类型转换器
	SetSqlArgTypeConvert(convert SqlArgTypeConvert)

	//表达式执行引擎
	ExpressionEngine() ExpressionEngine

	//设置表达式执行引擎
	SetExpressionEngine(engine ExpressionEngine)

	//sql构建器
	SqlBuilder() SqlBuilder

	//设置sql构建器
	SetSqlBuilder(builder SqlBuilder)

	//sql查询结果解析器
	SqlResultDecoder() SqlResultDecoder

	//设置sql查询结果解析器
	SetSqlResultDecoder(decoder SqlResultDecoder)
}

package GoMybatis

import (
	"database/sql"
)

type GoMybatisEngine struct {
	dbMap            map[string]*sql.DB //数据库map
	dataSourceRouter DataSourceRouter   //动态数据源路由器
	log              Log                //日志实现
	logEnable        bool               //是否允许日志输出

	sessionFactory *SessionFactory

	expressionTypeConvert ExpressionTypeConvert

	sqlArgTypeConvert SqlArgTypeConvert

	expressionEngine ExpressionEngine

	sqlBuilder SqlBuilder

	sqlResultDecoder SqlResultDecoder
}

func (it GoMybatisEngine) New() GoMybatisEngine {
	it.dbMap = make(map[string]*sql.DB)
	it.logEnable = true
	return it
}

func (it *GoMybatisEngine) WriteMapperPtr(ptr interface{}, xml []byte) {
	WriteMapperPtrByEngine(ptr, xml, it)
}

func (it *GoMybatisEngine) Name() string {
	return "GoMybatisEngine"
}

func (it *GoMybatisEngine) DataSourceRouter() DataSourceRouter {
	if it.dataSourceRouter == nil {
		var newRouter = GoMybatisDataSourceRouter{}.New(nil)
		DefaultGoMybatisEngine.SetDataSourceRouter(&newRouter)
	}
	return it.dataSourceRouter
}
func (it *GoMybatisEngine) SetDataSourceRouter(router DataSourceRouter) {
	for k, v := range it.dbMap {
		router.SetDB(k, v)
	}
	it.dataSourceRouter = router
}

func (it *GoMybatisEngine) DBMap() map[string]*sql.DB {
	return it.dbMap
}

func (it *GoMybatisEngine) NewSession(mapperName string) (Session, error) {
	var session, err = it.DataSourceRouter().Router(mapperName)
	return session, err
}

//获取日志实现类，是否启用日志
func (it *GoMybatisEngine) LogEnable() (Log, bool) {
	return it.log, it.logEnable
}

//设置日志实现类，是否启用日志
func (it *GoMybatisEngine) SetLogEnable(enable bool, log Log) {
	it.logEnable = enable
	it.log = log
}

//session工厂
func (it *GoMybatisEngine) SessionFactory() *SessionFactory {
	if it.sessionFactory == nil {
		var factory = SessionFactory{}.New(it)
		it.sessionFactory = &factory
	}
	return it.sessionFactory
}

//设置session工厂
func (it *GoMybatisEngine) SetSessionFactory(factory *SessionFactory) {
	it.sessionFactory = factory
}

//表达式数据类型转换器
func (it *GoMybatisEngine) ExpressionTypeConvert() ExpressionTypeConvert {
	if it.expressionTypeConvert == nil {
		it.expressionTypeConvert = GoMybatisExpressionTypeConvert{}
	}
	return it.expressionTypeConvert
}

//设置表达式数据类型转换器
func (it *GoMybatisEngine) SetExpressionTypeConvert(convert ExpressionTypeConvert) {
	it.expressionTypeConvert = convert
}

//sql类型转换器
func (it *GoMybatisEngine) SqlArgTypeConvert() SqlArgTypeConvert {
	if it.sqlArgTypeConvert == nil {
		it.sqlArgTypeConvert = GoMybatisSqlArgTypeConvert{}
	}
	return it.sqlArgTypeConvert
}

//设置sql类型转换器
func (it *GoMybatisEngine) SetSqlArgTypeConvert(convert SqlArgTypeConvert) {
	it.sqlArgTypeConvert = convert
}

//表达式执行引擎
func (it *GoMybatisEngine) ExpressionEngine() ExpressionEngine {
	if it.expressionEngine == nil {
		it.expressionEngine = &ExpressionEngineExpr{}
	}
	return it.expressionEngine
}

//设置表达式执行引擎
func (it *GoMybatisEngine) SetExpressionEngine(engine ExpressionEngine) {
	it.expressionEngine = engine
}

//sql构建器
func (it *GoMybatisEngine) SqlBuilder() SqlBuilder {
	if it.sqlBuilder == nil {
		var expressionEngineProxy = ExpressionEngineProxy{}.New(it.ExpressionEngine(), true)
		var log, enable = it.LogEnable()
		it.sqlBuilder = GoMybatisSqlBuilder{}.New(it.ExpressionTypeConvert(), it.SqlArgTypeConvert(), expressionEngineProxy, log, enable)
	}
	return it.sqlBuilder
}

//设置sql构建器
func (it *GoMybatisEngine) SetSqlBuilder(builder SqlBuilder) {
	it.sqlBuilder = builder
}

//sql查询结果解析器
func (it *GoMybatisEngine) SqlResultDecoder() SqlResultDecoder {
	if it.sqlResultDecoder == nil {
		it.sqlResultDecoder = GoMybatisSqlResultDecoder{}
	}
	return it.sqlResultDecoder
}

//设置sql查询结果解析器
func (it *GoMybatisEngine) SetSqlResultDecoder(decoder SqlResultDecoder) {
	it.sqlResultDecoder = decoder
}

//打开一个本地引擎
//driverName: 驱动名称例如"mysql", dataSourceName: string 数据库url
func Open(driverName, dataSourceName string) (SessionEngine, error) {
	if DefaultGoMybatisEngine == nil {
		var goMybatisEngine = GoMybatisEngine{}.New()
		DefaultGoMybatisEngine = SessionEngine(&goMybatisEngine)
	}
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	DefaultGoMybatisEngine.DBMap()[dataSourceName] = db
	return DefaultGoMybatisEngine, nil
}

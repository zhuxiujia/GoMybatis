package GoMybatis

import (
	"database/sql"
	"github.com/zhuxiujia/GoMybatis/utils"
)

type GoMybatisEngine struct {
	isInit           bool               //是否初始化
	dbMap            map[string]*sql.DB //数据库map（默认不为nil）
	dataSourceRouter DataSourceRouter   //动态数据源路由器
	log              Log                //日志实现
	logEnable        bool               //是否允许日志输出（默认开启）

	sessionFactory *SessionFactory

	sqlArgTypeConvert SqlArgTypeConvert

	expressionEngine ExpressionEngine

	sqlBuilder SqlBuilder

	sqlResultDecoder SqlResultDecoder

	templeteDecoder TempleteDecoder
}

func (it GoMybatisEngine) New() GoMybatisEngine {
	it.dbMap = make(map[string]*sql.DB)
	it.logEnable = true
	it.isInit = true

	if it.dataSourceRouter == nil {
		var newRouter = GoMybatisDataSourceRouter{}.New(nil)
		it.SetDataSourceRouter(&newRouter)
	}

	if it.logEnable == true && it.log == nil {
		it.log = &LogStandard{}
	}
	if it.sqlArgTypeConvert == nil {
		it.sqlArgTypeConvert = GoMybatisSqlArgTypeConvert{}
	}
	if it.expressionEngine == nil {
		it.expressionEngine = &ExpressionEngineExpr{}
	}
	if it.sqlResultDecoder == nil {
		it.sqlResultDecoder = GoMybatisSqlResultDecoder{}
	}
	if it.templeteDecoder == nil {
		it.SetTempleteDecoder(&GoMybatisTempleteDecoder{})
	}

	if it.sqlBuilder == nil {
		var expressionEngineProxy = ExpressionEngineProxy{}.New(it.ExpressionEngine(), true)
		var builder = GoMybatisSqlBuilder{}.New(it.SqlArgTypeConvert(), expressionEngineProxy, it.Log(), it.LogEnable())
		it.sqlBuilder = &builder
	}

	if it.sessionFactory == nil {
		var factory = SessionFactory{}.New(&it)
		it.sessionFactory = &factory
	}
	return it
}

func (it GoMybatisEngine) initCheck() {
	if it.isInit == false {
		panic(utils.NewError("GoMybatisEngine", "must call GoMybatisEngine{}.New() to init!"))
	}
}

func (it *GoMybatisEngine) WriteMapperPtr(ptr interface{}, xml []byte) {
	it.initCheck()
	WriteMapperPtrByEngine(ptr, xml, it)
}

func (it *GoMybatisEngine) Name() string {
	return "GoMybatisEngine"
}

func (it *GoMybatisEngine) DataSourceRouter() DataSourceRouter {
	it.initCheck()
	return it.dataSourceRouter
}
func (it *GoMybatisEngine) SetDataSourceRouter(router DataSourceRouter) {
	it.initCheck()
	for k, v := range it.dbMap {
		router.SetDB(k, v)
	}
	it.dataSourceRouter = router
}

func (it *GoMybatisEngine) DBMap() map[string]*sql.DB {
	it.initCheck()
	return it.dbMap
}

func (it *GoMybatisEngine) NewSession(mapperName string) (Session, error) {
	it.initCheck()
	var session, err = it.DataSourceRouter().Router(mapperName)
	return session, err
}

//获取日志实现类，是否启用日志
func (it *GoMybatisEngine) LogEnable() bool {
	it.initCheck()
	return it.logEnable
}

//设置日志实现类，是否启用日志
func (it *GoMybatisEngine) SetLogEnable(enable bool) {
	it.initCheck()
	it.logEnable = enable
	it.sqlBuilder.SetEnableLog(enable)
}

//获取日志实现类
func (it *GoMybatisEngine) Log() Log {
	it.initCheck()
	return it.log
}

//设置日志实现类
func (it *GoMybatisEngine) SetLog(log Log) {
	it.initCheck()
	it.log = log
}

//session工厂
func (it *GoMybatisEngine) SessionFactory() *SessionFactory {
	it.initCheck()
	return it.sessionFactory
}

//设置session工厂
func (it *GoMybatisEngine) SetSessionFactory(factory *SessionFactory) {
	it.initCheck()
	it.sessionFactory = factory
}

//sql类型转换器
func (it *GoMybatisEngine) SqlArgTypeConvert() SqlArgTypeConvert {
	it.initCheck()
	return it.sqlArgTypeConvert
}

//设置sql类型转换器
func (it *GoMybatisEngine) SetSqlArgTypeConvert(convert SqlArgTypeConvert) {
	it.initCheck()
	it.sqlArgTypeConvert = convert
}

//表达式执行引擎
func (it *GoMybatisEngine) ExpressionEngine() ExpressionEngine {
	it.initCheck()
	return it.expressionEngine
}

//设置表达式执行引擎
func (it *GoMybatisEngine) SetExpressionEngine(engine ExpressionEngine) {
	it.initCheck()
	it.expressionEngine = engine
	var proxy = it.sqlBuilder.ExpressionEngineProxy()
	proxy.SetExpressionEngine(it.expressionEngine)
}

//sql构建器
func (it *GoMybatisEngine) SqlBuilder() SqlBuilder {
	it.initCheck()
	return it.sqlBuilder
}

//设置sql构建器
func (it *GoMybatisEngine) SetSqlBuilder(builder SqlBuilder) {
	it.initCheck()
	it.sqlBuilder = builder
}

//sql查询结果解析器
func (it *GoMybatisEngine) SqlResultDecoder() SqlResultDecoder {
	it.initCheck()
	return it.sqlResultDecoder
}

//设置sql查询结果解析器
func (it *GoMybatisEngine) SetSqlResultDecoder(decoder SqlResultDecoder) {
	it.initCheck()
	it.sqlResultDecoder = decoder
}

//打开数据库
//driverName: 驱动名称例如"mysql", dataSourceName: string 数据库url
func (it *GoMybatisEngine) Open(driverName, dataSourceName string) error {
	it.initCheck()
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return err
	}
	it.DBMap()[dataSourceName] = db
	return nil
}

//模板解析器
func (it *GoMybatisEngine) TempleteDecoder() TempleteDecoder {
	return it.templeteDecoder
}

//设置模板解析器
func (it *GoMybatisEngine) SetTempleteDecoder(decoder TempleteDecoder) {
	it.templeteDecoder = decoder
}

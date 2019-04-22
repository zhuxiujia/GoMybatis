package GoMybatis

import (
	"database/sql"
	"github.com/zhuxiujia/GoMybatis/ast"
	"github.com/zhuxiujia/GoMybatis/engines"
	"github.com/zhuxiujia/GoMybatis/utils"
	"sync"
)

type GoMybatisEngine struct {
	mutex             sync.RWMutex          //读写锁
	isInit            bool                  //是否初始化
	dataSourceRouter  DataSourceRouter      //动态数据源路由器
	log               Log                   //日志实现
	logEnable         bool                  //是否允许日志输出（默认开启）
	sessionFactory    *SessionFactory       //session 工厂
	sqlArgTypeConvert ast.SqlArgTypeConvert //sql参数转换
	expressionEngine  ast.ExpressionEngine  //表达式解析引擎
	sqlBuilder        SqlBuilder            //sql 构建
	sqlResultDecoder  SqlResultDecoder      //sql查询结果解析引擎
	templeteDecoder   TempleteDecoder       //模板解析引擎
	callBackChain     []*CallBack           //回调链
	goroutineSessionMap *GoroutineSessionMap //map[协程id]Session
}

func (it GoMybatisEngine) New() GoMybatisEngine {
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
		it.expressionEngine = &engines.ExpressionEngineGoExpress{}
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
	if it.goroutineSessionMap==nil{
		var gr=GoroutineSessionMap{}.New()
		it.goroutineSessionMap=&gr
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
	it.dataSourceRouter = router
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
func (it *GoMybatisEngine) SqlArgTypeConvert() ast.SqlArgTypeConvert {
	it.initCheck()
	return it.sqlArgTypeConvert
}

//设置sql类型转换器
func (it *GoMybatisEngine) SetSqlArgTypeConvert(convert ast.SqlArgTypeConvert) {
	it.initCheck()
	it.sqlArgTypeConvert = convert
}

//表达式执行引擎
func (it *GoMybatisEngine) ExpressionEngine() ast.ExpressionEngine {
	it.initCheck()
	return it.expressionEngine
}

//设置表达式执行引擎
func (it *GoMybatisEngine) SetExpressionEngine(engine ast.ExpressionEngine) {
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
	it.dataSourceRouter.SetDB(dataSourceName, db)
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

//注册回调函数
func (it *GoMybatisEngine) RegisterCallBack(arg *CallBack) {
	it.mutex.Lock()
	if it.callBackChain == nil {
		it.callBackChain = make([]*CallBack, 0)
	}
	it.callBackChain = append(it.callBackChain, arg)
	it.mutex.Unlock()
}

func (it *GoMybatisEngine) CallBackChan() []*CallBack {
	return it.callBackChain
}

func (it *GoMybatisEngine) GoroutineSessionMap() *GoroutineSessionMap{
	return it.goroutineSessionMap
}
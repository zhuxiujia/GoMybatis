package GoMybatis

import (
	"database/sql"
)

type GoMybatisEngine struct {
	dbMap            map[string]*sql.DB
	dbMapLen         int
	dataSourceRouter DataSourceRouter
}

func (it GoMybatisEngine) New() GoMybatisEngine {
	it.dbMap = make(map[string]*sql.DB)
	return it
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

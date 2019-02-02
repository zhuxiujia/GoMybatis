package GoMybatis

import (
	"database/sql"
)

type GoMybatisEngine struct {
	SessionEngine
	dbMap            map[string]*sql.DB
	dbMapLen         int
	dataSourceRouter DataSourceRouter
}

func (it GoMybatisEngine) New() GoMybatisEngine {
	it.dbMap = make(map[string]*sql.DB)
	return it
}

func (it *GoMybatisEngine) DataSourceRouter() DataSourceRouter {
	return it.dataSourceRouter
}
func (it *GoMybatisEngine) SetDataSourceRouter(router DataSourceRouter) {
	it.dataSourceRouter = router
}

func (it *GoMybatisEngine) DBMap() map[string]*sql.DB {
	return it.dbMap
}

func (it *GoMybatisEngine) NewSession(mapperName string) (Session, error) {
	var session, err = it.dataSourceRouter.Router(mapperName)
	return session, err
}

var DefaultGoMybatisEngine SessionEngine

//打开一个本地引擎,driverName 驱动名称例如"mysql", dataSourceName string 数据库url, router DataSourceRouter 路由规则
func Open(driverName, dataSourceName string, router DataSourceRouter) (SessionEngine, error) {
	if DefaultGoMybatisEngine == nil {
		var goMybatisEngine = GoMybatisEngine{}.New()
		DefaultGoMybatisEngine = SessionEngine(&goMybatisEngine)
	}
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	DefaultGoMybatisEngine.DBMap()[dataSourceName] = db
	if router == nil {
		router = GoMybatisDataSourceRouter{}.New(db)
		DefaultGoMybatisEngine.SetDataSourceRouter(router)
	}
	return DefaultGoMybatisEngine, nil
}

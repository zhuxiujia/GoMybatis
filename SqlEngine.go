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
	//获取日志实现类，是否启用日志
	LogEnable()(Log,bool)
	//设置日志实现类，是否启用日志
	SetLogEnable(enable bool,log Log)
}

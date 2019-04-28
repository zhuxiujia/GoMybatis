package GoMybatis

import (
	"database/sql"
	"github.com/zhuxiujia/GoMybatis/utils"
)

//动态数据源路由
type GoMybatisDataSourceRouter struct {
	dbMap      map[string]*sql.DB
	driverMap  map[string]string
	routerFunc func(mapperName string) *string
}

//初始化路由，routerFunc为nil或者routerFunc返回nil，则框架自行选择第一个数据库作为数据源
func (it GoMybatisDataSourceRouter) New(routerFunc func(mapperName string) *string) GoMybatisDataSourceRouter {
	if routerFunc == nil {
		routerFunc = func(mapperName string) *string {
			return nil
		}
	}
	it.dbMap = make(map[string]*sql.DB)
	it.driverMap = make(map[string]string)
	it.routerFunc = routerFunc
	return it
}

func (it *GoMybatisDataSourceRouter) SetDB(driver string, url string, db *sql.DB) {
	it.dbMap[url] = db
	it.driverMap[url] = driver
}

func (it *GoMybatisDataSourceRouter) Router(mapperName string, engine SessionEngine) (Session, error) {
	var key *string
	var db *sql.DB

	if it.routerFunc != nil {
		key = it.routerFunc(mapperName)
	}

	if key != nil && *key != "" {
		db = it.dbMap[*key]
	} else {
		for k, v := range it.dbMap {
			if v != nil {
				db = v
				key = &k
				break
			}
		}
	}
	if db == nil {
		return nil, utils.NewError("GoMybatisDataSourceRouter", "router not find datasource opened ! do you forget invoke GoMybatis.GoMybatisEngine{}.New().Open(\"driverName\", Uri)?")
	}
	var url = ""
	if key != nil {
		url = *key
	}
	var local = LocalSession{}.New(it.driverMap[url], url, db, engine.Log())
	var session = Session(&local)
	return session, nil
}

func (it *GoMybatisDataSourceRouter) Name() string {
	return "GoMybatisDataSourceRouter"
}

package GoMybatis

import (
	"database/sql"
	"github.com/zhuxiujia/GoMybatis/utils"
)

type GoMybatisDataSourceRouter struct {
	dbMap      map[string]*sql.DB
	routerFunc func(mapperName string) *string
}

func (it GoMybatisDataSourceRouter) New(routerFunc func(mapperName string) *string) GoMybatisDataSourceRouter {
	if routerFunc == nil {
		routerFunc = func(mapperName string) *string {
			return nil
		}
	}
	it.dbMap = make(map[string]*sql.DB)
	it.routerFunc = routerFunc
	return it
}

func (it *GoMybatisDataSourceRouter) SetDB(url string, db *sql.DB) {
	it.dbMap[url] = db
}

func (it *GoMybatisDataSourceRouter) Router(mapperName string) (Session, error) {
	var key = it.routerFunc(mapperName)
	var db *sql.DB
	if key != nil && *key != "" {
		db = it.dbMap[*key]
	} else {
		for _, v := range it.dbMap {
			if v != nil {
				db = v
				break
			}
		}
		if db == nil {
			return nil, utils.NewError("GoMybatisDataSourceRouter", "router not find datasource!")
		}
	}
	var localSession = LocalSession{
		SessionId: utils.CreateUUID(),
		db:        db,
	}
	var session = Session(&localSession)
	return session, nil
}

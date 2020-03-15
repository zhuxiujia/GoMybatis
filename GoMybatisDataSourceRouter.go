package GoMybatis

import (
	"database/sql"
	"github.com/zhuxiujia/GoMybatis/utils"
)

//动态数据源路由
type GoMybatisDataSourceRouter struct {
	driverLinkDBMap  map[string]*sql.DB // map[driverLink]*DB
	driverLinkUrlMap map[string]string  // map[driverLink]Url
	routerFunc       func(mapperName string) *string
}

//初始化路由，routerFunc为nil或者routerFunc返回nil，则框架自行选择第一个数据库作为数据源
func (it GoMybatisDataSourceRouter) New(routerFunc func(mapperName string) *string) GoMybatisDataSourceRouter {
	if routerFunc == nil {
		routerFunc = func(mapperName string) *string {
			return nil
		}
	}
	it.driverLinkDBMap = make(map[string]*sql.DB)
	it.driverLinkUrlMap = make(map[string]string)
	it.routerFunc = routerFunc
	return it
}

func (it *GoMybatisDataSourceRouter) SetDB(driverType string, driverLink string, db *sql.DB) {
	it.driverLinkDBMap[driverLink] = db
	it.driverLinkUrlMap[driverLink] = driverType
}

func (it *GoMybatisDataSourceRouter) Router(mapperName string, engine SessionEngine) (Session, error) {
	var key *string
	var db *sql.DB

	if it.routerFunc != nil {
		key = it.routerFunc(mapperName)
	}

	if key != nil && *key != "" {
		db = it.driverLinkDBMap[*key]
	} else {
		for k, v := range it.driverLinkDBMap {
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
	var driverLink = ""
	if key != nil {
		driverLink = *key
	}
	var local = LocalSession{}.New(it.driverLinkUrlMap[driverLink], driverLink, db, engine.Log())
	var session = Session(&local)
	return session, nil
}

func (it *GoMybatisDataSourceRouter) Name() string {
	return "GoMybatisDataSourceRouter"
}

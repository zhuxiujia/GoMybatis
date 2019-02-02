package GoMybatis

import "database/sql"

//数据源路由接口
type DataSourceRouter interface {
	Router(mapperName string) (Session, error)
	SetDB(url string,db *sql.DB)
}

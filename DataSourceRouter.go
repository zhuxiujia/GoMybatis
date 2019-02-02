package GoMybatis

//数据源路由接口
type DataSourceRouter interface {
	Router(mapperName string) (Session, error)
}

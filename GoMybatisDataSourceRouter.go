package GoMybatis

import (
	"database/sql"
	"github.com/zhuxiujia/GoMybatis/utils"
)

type GoMybatisDataSourceRouter struct {
	db *sql.DB
}

func (it GoMybatisDataSourceRouter) New(db *sql.DB) GoMybatisDataSourceRouter {
	it.db = db
	return it
}

func (it GoMybatisDataSourceRouter) Router(mapperName string) (Session, error) {
	var localSession = LocalSession{
		SessionId: utils.CreateUUID(),
		db:        it.db,
	}
	var session = Session(&localSession)
	return session, nil
}

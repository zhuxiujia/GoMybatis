package GoMybatis

import (
	"database/sql"
	"github.com/zhuxiujia/GoMybatis/utils"
)

type GoMybatisEngine struct {
	SessionEngine
	DB *sql.DB
}

func (this GoMybatisEngine) NewSession() Session {
	uuid := utils.CreateUUID()
	var mysqlLocalSession = LocalSession{
		SessionId: uuid,
		db:        this.DB,
	}
	var session = Session(&mysqlLocalSession)
	return session
}

//打开一个本地引擎
func Open(driverName, dataSourceName string) (*SessionEngine, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	var mysqlEngine = GoMybatisEngine{
		DB: db,
	}
	var engine = SessionEngine(mysqlEngine)
	return &engine, nil
}

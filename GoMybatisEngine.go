package GoMybatis

import (
	"database/sql"
	"github.com/zhuxiujia/GoMybatis/lib/github.com/satori/go.uuid"
)

type GoMybatisEngine struct {
	SessionEngine
	DB *sql.DB
}

func (this GoMybatisEngine) NewSession() Session {
	uuids, _ := uuid.NewV4()
	var uuidstrig = uuids.String()
	var mysqlLocalSession = LocalSession{
		SessionId: uuidstrig,
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

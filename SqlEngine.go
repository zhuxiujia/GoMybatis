package GoMybatis

import "database/sql"

type Result struct {
	LastInsertId int64
	RowsAffected int64
}

type Session interface {
	SessionId() string
	DB() *sql.DB
	Query(sqlorArgs string) ([]map[string][]byte, error)
	Exec(sqlorArgs string) (Result, error)
	Rollback() error
	Commit() error
	Begin() error
	Close()
}

type SqlEngine interface {
	NewSession(sql string) *Session
	Open(driverName, dataSourceName string) (*SqlEngine, error)
}

package GoMybatis

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
	NewSession() Session
}

package GoMybatis

type Result struct {
	LastInsertId int64
	RowsAffected  int64
}

type SqlEngine interface {
	Query(sqlorArgs string) ([]map[string][]byte,error)
	Exec(sqlorArgs string) (Result, error)
}


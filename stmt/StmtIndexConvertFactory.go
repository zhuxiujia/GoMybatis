package stmt

import (
	"fmt"
	"sync"
)

// build a stmt convert
func BuildStmtConvert(driverType string) (StmtIndexConvert, error) {
	switch driverType {
	case "mysql", "mymysql", "mssql", "sqlite3":
		return &MysqlStmtIndexConvertImpl{}, nil
	case "postgres":
		return &PostgreStmtIndexConvertImpl{
			sync.RWMutex{},
			0,
		}, nil
	case "oci8":
		return &OracleStmtIndexConvertImpl{sync.RWMutex{},
			0}, nil
	default:
		panic(fmt.Sprint("[GoMybatis] un support dbName:", driverType, " only support: ", "mysql,", "mymysql,", "mssql,", "sqlite3,", "postgres,", "oci8"))
	}
}

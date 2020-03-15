package stmt

import "fmt"

// build a stmt convert
func BuildStmtConvert(dbName string) (StmtIndexConvert, error) {
	switch dbName {
	case "mysql", "mymysql", "mssql", "sqlite3":
		return &MysqlStmtIndexConvertImpl{}, nil
	case "postgres":
		return &MysqlStmtIndexConvertImpl{}, nil
	case "oci8":
		return &OracleStmtIndexConvertImpl{}, nil
	default:
		panic(fmt.Sprint("[GoMybatis] un support dbName:", dbName, " only support: ", "mysql,", "mymysql,", "mssql,", "sqlite3,", "postgres,", "oci8"))
	}
}

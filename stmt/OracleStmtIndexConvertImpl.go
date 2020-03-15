package stmt

import "fmt"

type OracleStmtIndexConvertImpl struct {
}

func (it *OracleStmtIndexConvertImpl) Convert(index int) string {
	return fmt.Sprint(" :val", index+1, " ")
}

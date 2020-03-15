package stmt

import "fmt"

type PostgreStmtIndexConvertImpl struct {
}

func (it *PostgreStmtIndexConvertImpl) Convert(index int) string {
	return fmt.Sprint(" $", index+1, " ")
}

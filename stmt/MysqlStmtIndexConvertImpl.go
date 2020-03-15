package stmt

type MysqlStmtIndexConvertImpl struct {
}

func (it *MysqlStmtIndexConvertImpl) Convert(index int) string {
	return " ? "
}

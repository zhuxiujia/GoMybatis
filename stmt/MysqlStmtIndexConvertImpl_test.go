package stmt

import "testing"

func TestMysqlStmtIndexConvertImpl_Convert(t *testing.T) {
	var convert = MysqlStmtIndexConvertImpl{}
	if " ? " != convert.Convert() {
		panic("TestMysqlStmtIndexConvertImpl_Convert fail")
	}
}

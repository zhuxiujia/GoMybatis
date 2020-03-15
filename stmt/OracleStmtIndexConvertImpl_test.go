package stmt

import "testing"

func TestOracleStmtIndexConvertImpl_Convert(t *testing.T) {
	var convert = OracleStmtIndexConvertImpl{}
	if " :val1 " != convert.Convert(0) {
		panic("TestOracleStmtIndexConvertImpl_Convert fail")
	}
}

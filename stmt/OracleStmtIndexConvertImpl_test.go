package stmt

import "testing"

func TestOracleStmtIndexConvertImpl_Convert(t *testing.T) {
	var convert = &OracleStmtIndexConvertImpl{}
	convert.Inc()
	if " :val1 " != convert.Convert() {
		panic("TestOracleStmtIndexConvertImpl_Convert fail")
	}
}

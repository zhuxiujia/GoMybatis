package stmt

import "testing"

func TestPostgreStmtIndexConvertImpl_Convert(t *testing.T) {
	var convert = PostgreStmtIndexConvertImpl{}
	if " $1 " != convert.Convert(0) {
		panic("TestPostgreStmtIndexConvertImpl_Convert fail")
	}
}

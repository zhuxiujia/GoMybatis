package stmt

import "testing"

func TestPostgreStmtIndexConvertImpl_Convert(t *testing.T) {
	var convert = &PostgreStmtIndexConvertImpl{}
	convert.Inc()
	convert.Inc()
	if " $2 " != convert.Convert() {
		panic("TestPostgreStmtIndexConvertImpl_Convert fail")
	}
}

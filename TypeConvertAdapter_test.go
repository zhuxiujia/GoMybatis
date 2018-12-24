package GoMybatis

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_Adapter(t *testing.T) {
	var a = true
	var convertResult = GoMybatisExpressionTypeConvert{}.Convert(SqlArg{
		Value: a,
		Type:  reflect.TypeOf(a),
	})
	if convertResult != true{
		t.Fatal(`Test_Adapter fail convertResult != true`)
	}
	fmt.Println(convertResult)
}

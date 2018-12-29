package GoMybatis

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_ExpressionTypeConvert(t *testing.T) {
	var a = true
	var convertResult = GoMybatisExpressionTypeConvert{}.Convert(SqlArg{
		Value: a,
		Type:  reflect.TypeOf(a),
	})
	if convertResult != true {
		t.Fatal(`Test_Adapter fail convertResult != true`)
	}
	fmt.Println(convertResult)
}

func Test_SqlArgTypeConvert(t *testing.T) {
	var a = true
	var convertResult = GoMybatisSqlArgTypeConvert{}.Convert(SqlArg{
		Value: a,
		Type:  reflect.TypeOf(a),
	})
	if convertResult == "" {
		t.Fatal(`Test_Adapter fail convertResult != true`)
	}
	fmt.Println(convertResult)
}

package GoMybatis

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_Adapter(t *testing.T) {
	var a bool
	var s = GoMybatisExpressionTypeConvert{}.Convert(SqlArg{
		Value: a,
		Type:  reflect.TypeOf(a),
	})
	fmt.Println(s)
}

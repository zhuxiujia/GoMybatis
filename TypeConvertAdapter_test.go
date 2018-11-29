package GoMybatis

import (
	"fmt"
	"testing"
	"reflect"
)

func Test_Adapter(t *testing.T) {
	var a bool
	var s = GoMybatisExpressionTypeConvert{}.Convert(SqlArg{
		Value:a,
		Type:reflect.TypeOf(a),
	})
	fmt.Println(s)
}

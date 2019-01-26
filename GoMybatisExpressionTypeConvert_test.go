package GoMybatis

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func Test_ExpressionTypeConvert(t *testing.T) {
	var a = true
	var convertResult = GoMybatisExpressionTypeConvert{}.Convert(SqlArg{
		Value: a,
		Type:  reflect.TypeOf(a),
	})
	if convertResult != true {
		t.Fatal(`Test_ExpressionTypeConvert fail`)
	}
	fmt.Println(convertResult)
	convertResult = GoMybatisExpressionTypeConvert{}.Convert(SqlArg{
		Value: 1,
		Type:  reflect.TypeOf(1),
	})
	if convertResult == nil {
		t.Fatal(`Test_ExpressionTypeConvert fail`)
	}
	fmt.Println(convertResult)
	convertResult = GoMybatisExpressionTypeConvert{}.Convert(SqlArg{
		Value: time.Now(),
		Type:  reflect.TypeOf(time.Now()),
	})
	if convertResult == nil {
		t.Fatal(`Test_ExpressionTypeConvert fail`)
	}
	fmt.Println(convertResult)
	convertResult = GoMybatisExpressionTypeConvert{}.Convert(SqlArg{
		Value: "string",
		Type:  reflect.TypeOf("string"),
	})
	if convertResult == nil {
		t.Fatal(`Test_ExpressionTypeConvert fail`)
	}
	fmt.Println(convertResult)
}

func BenchmarkGoMybatisExpressionTypeConvert_Convert(b *testing.B) {
	b.StopTimer()
	var convertType = reflect.TypeOf(1)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		var convertResult = GoMybatisExpressionTypeConvert{}.Convert(SqlArg{
			Value: 1,
			Type:  convertType,
		})
		if convertResult == nil {
			b.Fatal("convert fail!")
		}
	}
}

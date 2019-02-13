package GoMybatis

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func Test_ExpressionTypeConvert(t *testing.T) {
	var a = true
	var convertResult = GoMybatisExpressionTypeConvert{}.Convert( a, reflect.TypeOf(a))
	if convertResult != true {
		t.Fatal(`Test_ExpressionTypeConvert fail`)
	}
	fmt.Println(convertResult)
	convertResult = GoMybatisExpressionTypeConvert{}.Convert( 1, reflect.TypeOf(1))
	if convertResult == nil {
		t.Fatal(`Test_ExpressionTypeConvert fail`)
	}
	fmt.Println(convertResult)
	convertResult = GoMybatisExpressionTypeConvert{}.Convert(time.Now(), reflect.TypeOf(time.Now()))
	if convertResult == nil {
		t.Fatal(`Test_ExpressionTypeConvert fail`)
	}
	fmt.Println(convertResult)
	convertResult = GoMybatisExpressionTypeConvert{}.Convert("string", reflect.TypeOf("string"))
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
		var convertResult = GoMybatisExpressionTypeConvert{}.Convert(1, convertType)
		if convertResult == nil {
			b.Fatal("convert fail!")
		}
	}
}

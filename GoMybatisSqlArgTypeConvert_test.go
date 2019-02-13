package GoMybatis

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func Test_SqlArgTypeConvert(t *testing.T) {
	var a = true
	var convertResult = GoMybatisSqlArgTypeConvert{}.Convert( a,reflect.TypeOf(a))
	if convertResult == "" {
		t.Fatal(`Test_Adapter fail convertResult != true`)
	}
	fmt.Println(convertResult)
	convertResult = GoMybatisSqlArgTypeConvert{}.Convert( 1, reflect.TypeOf(1))
	if convertResult == "" {
		t.Fatal(`Test_Adapter fail convertResult != 1`)
	}
	fmt.Println(convertResult)
	convertResult = GoMybatisSqlArgTypeConvert{}.Convert(time.Now(),reflect.TypeOf(time.Now()))
	if convertResult == "" {
		t.Fatal(`Test_Adapter fail convertResult != time.Time`)
	}
	fmt.Println(convertResult)
	convertResult = GoMybatisSqlArgTypeConvert{}.Convert("string", reflect.TypeOf("string"))
	if convertResult == "" {
		t.Fatal(`Test_Adapter fail convertResult != string`)
	}
	fmt.Println(convertResult)
}

func Test_SqlArgTypeConvert_NoType(t *testing.T) {
	var a = true
	var convertResult = GoMybatisSqlArgTypeConvert{}.Convert( a,nil)
	if convertResult == "" {
		t.Fatal(`Test_Adapter fail convertResult != true`)
	}
	fmt.Println(convertResult)
	convertResult = GoMybatisSqlArgTypeConvert{}.Convert( 1, nil)
	if convertResult == "" {
		t.Fatal(`Test_Adapter fail convertResult != 1`)
	}
	fmt.Println(convertResult)
	convertResult = GoMybatisSqlArgTypeConvert{}.Convert(time.Now(),nil)
	if convertResult == "" {
		t.Fatal(`Test_Adapter fail convertResult != time.Time`)
	}
	fmt.Println(convertResult)
	convertResult = GoMybatisSqlArgTypeConvert{}.Convert("string", nil)
	if convertResult == "" {
		t.Fatal(`Test_Adapter fail convertResult != string`)
	}
	fmt.Println(convertResult)
}

func BenchmarkGoMybatisSqlArgTypeConvert_Convert(b *testing.B) {
	b.StopTimer()
	var convertType = reflect.TypeOf(1)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		var convertResult = GoMybatisSqlArgTypeConvert{}.Convert(1, convertType)
		if convertResult == "" {
			b.Fatal("convert fail!")
		}
	}
}

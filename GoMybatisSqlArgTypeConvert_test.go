package GoMybatis

import (
	"fmt"
	"testing"
	"time"
)

func Test_SqlArgTypeConvert(t *testing.T) {
	var a = true
	var convertResult = GoMybatisSqlArgTypeConvert{}.Convert(a)
	if convertResult != "true" {
		t.Fatal(`Test_Adapter fail convertResult != true`)
	}
	fmt.Println(convertResult)
	convertResult = GoMybatisSqlArgTypeConvert{}.Convert(1)
	if convertResult != "1" {
		t.Fatal(`Test_Adapter fail convertResult != 1`)
	}
	fmt.Println(convertResult)
	var now = time.Now()
	convertResult = GoMybatisSqlArgTypeConvert{}.Convert(now)
	if convertResult != "'"+now.Format(Adapter_FormateDate)+"'" {
		t.Fatal(`Test_Adapter fail convertResult != 2019-05-10 11:09:01`)
	}
	fmt.Println(convertResult)
	convertResult = GoMybatisSqlArgTypeConvert{}.Convert("string")
	if convertResult != "'string'" {
		t.Fatal(`Test_Adapter fail convertResult != string`)
	}
	fmt.Println(convertResult)
}

func Test_SqlArgTypeConvert_NoType(t *testing.T) {
	var a = true
	var convertResult = GoMybatisSqlArgTypeConvert{}.Convert(a)
	if convertResult == "" {
		t.Fatal(`Test_Adapter fail convertResult != true`)
	}
	fmt.Println(convertResult)
	convertResult = GoMybatisSqlArgTypeConvert{}.Convert(1)
	if convertResult == "" {
		t.Fatal(`Test_Adapter fail convertResult != 1`)
	}
	fmt.Println(convertResult)
	convertResult = GoMybatisSqlArgTypeConvert{}.Convert(time.Now())
	if convertResult == "" {
		t.Fatal(`Test_Adapter fail convertResult != time.Time`)
	}
	fmt.Println(convertResult)
	convertResult = GoMybatisSqlArgTypeConvert{}.Convert("string")
	if convertResult == "" {
		t.Fatal(`Test_Adapter fail convertResult != string`)
	}
	fmt.Println(convertResult)
}

func BenchmarkGoMybatisSqlArgTypeConvert_Convert(b *testing.B) {
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		var convertResult = GoMybatisSqlArgTypeConvert{}.Convert(1)
		if convertResult == "" {
			b.Fatal("convert fail!")
		}
	}
}

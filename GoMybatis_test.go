package GoMybatis

import (
	"github.com/zhuxiujia/GoMybatis/example"
	"reflect"
	"testing"
	"time"
)

func TestScanStructArgFields(ff *testing.T) {
	var act = example.Activity{
		Id:         "123",
		Uuid:       "uu",
		Name:       "test",
		PcLink:     "pc",
		H5Link:     "h5",
		Remark:     "remark",
		Version:    0,
		CreateTime: time.Now(),
		DeleteFlag: 1,
	}
	scanStructArgFields(reflect.ValueOf(act), nil)
	var t = reflect.TypeOf(act)
	for i := 0; i < t.NumField(); i++ {
		var typeValue = t.Field(i)
		var jsonKey = typeValue.Tag.Get(`json`)
		println(jsonKey)
	}
}

func BenchmarkScanStructArgFields(b *testing.B) {
	b.StopTimer()
	var act = example.Activity{
		Id:         "123",
		Uuid:       "uu",
		Name:       "test",
		PcLink:     "pc",
		H5Link:     "h5",
		Remark:     "remark",
		Version:    0,
		CreateTime: time.Now(),
		DeleteFlag: 1,
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		scanStructArgFields(reflect.ValueOf(act), nil)
	}
}

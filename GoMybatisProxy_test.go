package GoMybatis

import (
	"fmt"
	"reflect"
	"testing"
)

type TestMapper struct {
	SelectByIds func(id string) (string, error) `mapperParams:"ids"`
}

func TestUseMapperValue(t *testing.T) {
	var test = TestMapper{}
	UseMapperValue(reflect.ValueOf(&test), func(method string, args []reflect.Value, tagArgs []TagArg) []reflect.Value {
		if method == "" {
			t.Fatal("UseMapper() method == ''")
		}
		if len(args) <= 0 {
			t.Fatal("UseMapper() args len = 0")
		}
		if len(tagArgs) <= 0 {
			t.Fatal("UseMapper() tagArgs len = 0")
		}
		var e error
		var returns = make([]reflect.Value, 0)
		returns = append(returns, reflect.ValueOf("yes return string="+args[0].Interface().(string)))
		returns = append(returns, reflect.Zero(reflect.TypeOf(&e).Elem()))
		return returns
	})

	var result, e = test.SelectByIds("1234")
	fmt.Println(result, e)
	if e != nil {
		t.Fatal(e)
	}
}

func TestUseMapper(t *testing.T) {
	var test = TestMapper{}
	UseMapper(&test, func(method string, args []reflect.Value, tagArgs []TagArg) []reflect.Value {
		if method == "" {
			t.Fatal("UseMapper() method == ''")
		}
		if len(args) <= 0 {
			t.Fatal("UseMapper() args len = 0")
		}
		if len(tagArgs) <= 0 {
			t.Fatal("UseMapper() tagArgs len = 0")
		}
		var e error
		var returns = make([]reflect.Value, 0)
		returns = append(returns, reflect.ValueOf("yes return string="+args[0].Interface().(string)))
		returns = append(returns, reflect.Zero(reflect.TypeOf(&e).Elem()))
		return returns
	})

	var result, e = test.SelectByIds("1234")
	fmt.Println(result, e)
	if e != nil {
		t.Fatal(e)
	}
}

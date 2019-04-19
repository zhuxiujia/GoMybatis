package example

import (
	"github.com/zhuxiujia/GoMybatis"
	"reflect"
	"testing"
)

type Service struct {
	FindName func() (string, error)
}

func TestService(t *testing.T) {
	var s = Service{
		FindName: func() (string, error) {
			return "TestService", nil
		},
	}
	AopProxyService(&s)
	var r, _ = s.FindName()
	println("result:", r)
}

func AopProxyService(service interface{}) {
	GoMybatis.UseMapper(service, func(funcField reflect.StructField, field reflect.Value) func(arg GoMybatis.ProxyArg) []reflect.Value {
		//拷贝老方法，否则会循环调用导致栈溢出
		var oldFunc = reflect.ValueOf(field.Interface())
		var fn = func(arg GoMybatis.ProxyArg) []reflect.Value {
			//befor

			var oldFuncResults = oldFunc.Call(arg.Args)
			//after

			return oldFuncResults
		}
		return fn
	})
}

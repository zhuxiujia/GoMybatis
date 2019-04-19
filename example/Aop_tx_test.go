package example

import (
	"github.com/zhuxiujia/GoMybatis"
	"reflect"
	"testing"
)

type Service struct {
	FindName func() (string, error)

	SayHello func() (string, error)
}

func TestService(t *testing.T) {
	var s Service
	s = Service{
		FindName: func() (string, error) {
			println("TestService")
			s.SayHello()
			return "TestService", nil
		},
		SayHello: func() (s string, e error) {
			println("hello")
			return "hello", nil
		},
	}
	AopProxyService(&s)
	s.FindName()
}

func AopProxyService(service interface{}) {
	GoMybatis.AopProxy(service, func(funcField reflect.StructField, field reflect.Value) func(arg GoMybatis.ProxyArg) []reflect.Value {
		//拷贝老方法，否则会循环调用导致栈溢出
		var oldFunc = reflect.ValueOf(field.Interface())
		var fn = func(arg GoMybatis.ProxyArg) []reflect.Value {
			println("start:" + funcField.Name)
			var oldFuncResults = oldFunc.Call(arg.Args)
			println("end:" + funcField.Name)
			return oldFuncResults
		}
		return fn
	})
}

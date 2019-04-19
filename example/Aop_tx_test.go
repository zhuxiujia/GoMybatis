package example

import (
	"github.com/zhuxiujia/GoMybatis"
	"reflect"
	"testing"
)

type Service struct {
	FindName func() error
	SayHello func() error
}

func TestService(t *testing.T) {
	var s = Service{}
	s.FindName = func() error {
		println("TestService")
		s.SayHello()
		return nil
	}
	s.SayHello = func() error {
		println("hello")
		return nil
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

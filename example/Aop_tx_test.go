package example

import (
	"github.com/zhuxiujia/GoMybatis"
	"reflect"
	"testing"
)

type Service struct {
	FindName func(it *Service) error `transaction:"PROPAGATION_REQUIRED"`
	SayHello func(it *Service) error
}

func TestService(t *testing.T) {
	var s = Service{
		FindName: func(it *Service) error {
			println("TestService")
			it.SayHello(it)
			return nil
		},
		SayHello: func(it *Service) error {
			println("hello")
			return nil
		},
	}
	AopProxyService(&s)
	s.FindName(&s)
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

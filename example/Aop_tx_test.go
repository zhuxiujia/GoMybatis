package example

import (
	"github.com/zhuxiujia/GoMybatis"
	"github.com/zhuxiujia/GoMybatis/tx"
	"reflect"
	"testing"
)

type TestService struct {
	FindName func() error `transaction:"PROPAGATION_REQUIRED"`
	SayHello func() error
}

func TestTestService(t *testing.T) {
	var it TestService
	it = TestService{
		FindName: func() error {
			println("TestService")
			it.SayHello()
			return nil
		},
		SayHello: func() error {
			println("hello")
			return nil
		},
	}
	AopProxyService(&it)
	it.FindName()
}

func AopProxyService(service interface{}) {
	//调用方法栈
	var txStack = tx.StructField{}.New()
	GoMybatis.Proxy(service, func(funcField reflect.StructField, field reflect.Value) func(arg GoMybatis.ProxyArg) []reflect.Value {
		//拷贝老方法，否则会循环调用导致栈溢出
		var nativeImplFunc = reflect.ValueOf(field.Interface())
		var fn = func(arg GoMybatis.ProxyArg) []reflect.Value {
			txStack.Push(funcField)
			var nativeImplResult = nativeImplFunc.Call(arg.Args)
			txStack.Pop()
			if txStack.Len() == 0 {
				//todo rollback
			}
			return nativeImplResult
		}
		return fn
	})
}

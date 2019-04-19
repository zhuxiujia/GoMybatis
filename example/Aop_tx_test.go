package example

import (
	"github.com/zhuxiujia/GoMybatis"
	"github.com/zhuxiujia/GoMybatis/tx"
	"reflect"
	"testing"
)

type Service struct {
	FindName func() error `transaction:"PROPAGATION_REQUIRED"`
	SayHello func() error
}

func TestService(t *testing.T) {
	var it Service
	it = Service{
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
	var txStack = tx.TxStack{}.New()
	GoMybatis.AopProxy(service, func(funcField reflect.StructField, field reflect.Value) func(arg GoMybatis.ProxyArg) []reflect.Value {
		//拷贝老方法，否则会循环调用导致栈溢出
		var oldFunc = reflect.ValueOf(field.Interface())
		var fn = func(arg GoMybatis.ProxyArg) []reflect.Value {
			txStack.Push(funcField)
			var oldFuncResults = oldFunc.Call(arg.Args)
			txStack.Pop()
			if txStack.Len() == 0 {
				//todo rollback
			}
			return oldFuncResults
		}
		return fn
	})
}

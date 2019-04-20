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
	var engine = GoMybatis.GoMybatisEngine{}
	var beanType = reflect.TypeOf(service).Elem()
	var beanName = beanType.PkgPath() + beanType.Name()
	var session GoMybatis.Session
	var txStack = tx.StructField{}.New()
	GoMybatis.Proxy(service, func(funcField reflect.StructField, field reflect.Value) func(arg GoMybatis.ProxyArg) []reflect.Value {
		//拷贝老方法，否则会循环调用导致栈溢出
		var nativeImplFunc = reflect.ValueOf(field.Interface())
		var fn = func(arg GoMybatis.ProxyArg) []reflect.Value {
			txStack.Push(funcField)
			if txStack.Len() == 1 {
				//PROPAGATION_REQUIRED
				if session == nil {
					//todo newSession is use service bean name?
					session, _ = engine.NewSession(beanName)
					session.Begin()
				}
			}
			defer func() {
				err := recover()
				if err != nil {
					session.Rollback()
				}
			}()
			var nativeImplResult = nativeImplFunc.Call(arg.Args)
			txStack.Pop()
			if txStack.Len() == 0 {
				if haveError(nativeImplResult) {
					session.Commit()
				} else {
					session.Rollback()
				}
			}
			return nativeImplResult
		}
		return fn
	})
}

func haveError(v []reflect.Value) bool {
	if v == nil || len(v) == 0 {
		return false
	}
	for _, item := range v {
		if item.Kind() == reflect.Interface {
			if item.String() == "error" {
				if !item.IsNil() {
					return true
				}
			}
		}
	}
	return false
}

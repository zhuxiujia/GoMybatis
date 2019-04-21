package GoMybatis

import (
	"github.com/zhuxiujia/GoMybatis/tx"
	"reflect"
)

//使用AOP切面 代理目标服务，如果服务painc()它的事务会回滚
func AopProxyService(service reflect.Value, engine *GoMybatisEngine) {
	//调用方法栈
	var beanType = service.Type().Elem()
	var beanName = beanType.PkgPath() + beanType.Name()
	var session Session
	var txStack = tx.StructField{}.New()
	ProxyValue(service, func(funcField reflect.StructField, field reflect.Value) func(arg ProxyArg) []reflect.Value {
		//init data
		var propagation = tx.PROPAGATION_NEVER
		var nativeImplFunc = reflect.ValueOf(field.Interface())
		var txTag, haveTx = funcField.Tag.Lookup("tx")
		var rollbackTag = funcField.Tag.Get("rollback")
		if haveTx {
			propagation = tx.NewPropagation(txTag)
		}
		var fn = func(arg ProxyArg) []reflect.Value {
			txStack.Push(funcField)
			if txStack.Len() == 1 {
				if propagation == tx.PROPAGATION_NEVER{

				} else if propagation == tx.PROPAGATION_REQUIRED {
					//PROPAGATION_REQUIRED
					if session == nil {
						//todo newSession is use service bean name?
						var err error
						session, err = engine.NewSession(beanName)
						if err != nil {
							panic(err)
						}
						err = session.Begin()
						if err != nil {
							panic(err)
						}
						println("Begin in session:",session.Id())
					}
				}
			}
			var nativeImplResult = doNativeMethod(arg, nativeImplFunc, session)
			txStack.Pop()
			if txStack.Len() == 0 && session != nil {
				if !haveRollBackType(nativeImplResult, rollbackTag) {
					var err = session.Commit()
					if err != nil {
						panic(err)
					}
					println("Commit in session:",session.Id())
				} else {
					var err = session.Rollback()
					if err != nil {
						panic(err)
					}
					println("Rollback in session:",session.Id())
				}
			}
			return nativeImplResult
		}
		return fn
	})
}

func doNativeMethod(arg ProxyArg, nativeImplFunc reflect.Value, session Session) []reflect.Value {
	defer func() {
		err := recover()
		if err != nil {
			var err = session.Rollback()
			if err != nil {
				panic(err)
			}
		}
	}()
	return nativeImplFunc.Call(arg.Args)

}

func haveRollBackType(v []reflect.Value, typeString string) bool {
	if v == nil || len(v) == 0 || typeString == ""{
		return false
	}
	for _, item := range v {
		if item.Kind() == reflect.Interface {
			if item.String() == typeString {
				if !item.IsNil() {
					return true
				}
			}
		}
	}
	return false
}

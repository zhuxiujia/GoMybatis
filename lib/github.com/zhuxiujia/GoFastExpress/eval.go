package GoFastExpress

import (
	"fmt"
	"github.com/kataras/iris/core/errors"
	"reflect"
	"strings"
)

func EvalMap(arg map[string]interface{}) {

}

//取值
func EvalTake(operator Operator, arg interface{}) (interface{}, error) {
	if arg == nil {
		return nil, nil
	}
	if operator == "" {
		return arg, nil
	}
	var av = reflect.ValueOf(arg)
	if av.Kind() == reflect.Map {
		var m = arg.(map[string]interface{})
		var result = m[operator]
		if result != nil {
			return result, nil
		}
		if strings.Index(operator, ".") != -1 {
			var item []byte
			for index, v := range operator {
				if v == 46 {
					item = []byte(operator)[0:index]
					break
				}
			}
			result = m[string(item)]
			var otheritem = string([]byte(operator)[len(item)+1 : len(operator)])
			return getObj(otheritem, reflect.ValueOf(result), result)
		}
		return nil, nil
	} else {
		return getObj(operator, av, arg)
	}
	return arg, nil
}

func getObj(operator Operator, av reflect.Value, arg interface{}) (interface{}, error) {
	if av.Kind() == reflect.Ptr {
		av = GetDeepPtr(av)
	}
	var v = av.FieldByName(operator)
	if v.IsValid() && v.CanInterface() {
		return v.Interface(), nil
	} else {
		return nil, nil
	}
	return arg, nil
}

func Eval(operator Operator, a interface{}, b interface{}) (interface{}, error) {
	var av = reflect.ValueOf(a)
	var bv = reflect.ValueOf(b)

	switch operator {
	case And:
		if a == nil || b == nil {
			//equal nil
			return nil, errors.New("eval fail,value can not be nil")
		}
		a, av = GetDeepValue(av, a)
		b, bv = GetDeepValue(bv, b)
		var ab = a.(bool)
		var bb = b.(bool)
		return ab == true && bb == true, nil
	case Or:
		if a == nil || b == nil {
			//equal nil
			return nil, errors.New("eval fail,value can not be nil")
		}
		a, av = GetDeepValue(av, a)
		b, bv = GetDeepValue(bv, b)
		var ab = a.(bool)
		var bb = b.(bool)
		return ab == true || bb == true, nil
	case Equal, MoreEqual, More, Less, LessEqual:
		//a kind == b kind
		return DoEqualAction(operator, a, b, av, bv)
	case UnEqual:
		//a kind == b kind
		var r, e = DoEqualAction(operator, a, b, av, bv)
		if e != nil {
			return nil, e
		}
		return !r, nil
	case Add, Reduce, Ride, Divide:
		return DoCalculationAction(operator, a, b, av, bv)
	}
	return nil, errors.New("find not support operator :" + operator)
}

func DoEqualAction(operator Operator, a interface{}, b interface{}, av reflect.Value, bv reflect.Value) (bool, error) {
	switch operator {
	case UnEqual:
		fallthrough
	case Equal:
		if av.Kind() == reflect.Ptr && av.IsNil() == true {
			a = nil
		}
		if bv.Kind() == reflect.Ptr && bv.IsNil() == true {
			b = nil
		}
		if a == nil || b == nil {
			//equal nil
			if a != nil || b != nil {
				return false, nil
			}
			if a == nil && b == nil {
				return true, nil
			}
		}
		if av.Kind() == reflect.Ptr {
			a, av = GetDeepValue(av, a)
		}
		if bv.Kind() == reflect.Ptr {
			b, bv = GetDeepValue(bv, b)
		}
		if av.Kind() == reflect.Float64 && bv.Kind() == reflect.Float64 {
			return a.(float64) == b.(float64), nil
		}
		if av.Kind() == reflect.Float32 && bv.Kind() == reflect.Float32 {
			return a.(float32) == b.(float32), nil
		}
		if av.Kind() == reflect.Int && bv.Kind() == reflect.Int {
			return a.(int) == b.(int), nil
		}
		if av.Kind() == reflect.Int8 && bv.Kind() == reflect.Int8 {
			return a.(int8) == b.(int8), nil
		}
		if av.Kind() == reflect.Int16 && bv.Kind() == reflect.Int16 {
			return a.(int16) == b.(int16), nil
		}
		if av.Kind() == reflect.Int32 && bv.Kind() == reflect.Int32 {
			return a.(int32) == b.(int32), nil
		}
		if av.Kind() == reflect.Int64 && bv.Kind() == reflect.Int64 {
			return a.(int64) == b.(int64), nil
		}
		if av.Kind() == reflect.Bool && bv.Kind() == reflect.Bool {
			return a.(bool) == b.(bool), nil
		}
		if av.Kind() == reflect.String && bv.Kind() == reflect.String {
			return a.(string) == b.(string), nil
		}
		if av.Kind() == reflect.Struct && bv.Kind() == reflect.String {
			return fmt.Sprint(a) == b.(string), nil
		}
		if bv.Kind() == reflect.Struct && av.Kind() == reflect.String {
			return fmt.Sprint(b) == a.(string), nil
		}
		a = toNumberType(av)
		b = toNumberType(bv)
		return a.(float64) == b.(float64), nil
	case Less:
		if a == nil || b == nil {
			return false, errors.New("can not parser '<' , arg have nil object!")
		}
		a, av = GetDeepValue(av, a)
		b, bv = GetDeepValue(bv, b)
		a = toNumberType(av)
		b = toNumberType(bv)
		return a.(float64) < b.(float64), nil
	case More:
		if a == nil || b == nil {
			return false, errors.New("can not parser '>' , arg have nil object!")
		}
		a, av = GetDeepValue(av, a)
		b, bv = GetDeepValue(bv, b)
		a = toNumberType(av)
		b = toNumberType(bv)
		return a.(float64) > b.(float64), nil
	case MoreEqual:
		if a == nil || b == nil {
			return false, errors.New("can not parser '>=' , arg have nil object!")
		}
		a, av = GetDeepValue(av, a)
		b, bv = GetDeepValue(bv, b)
		a = toNumberType(av)
		b = toNumberType(bv)
		return a.(float64) >= b.(float64), nil
	case LessEqual:
		if a == nil || b == nil {
			return false, errors.New("can not parser '<=' , arg have nil object!")
		}
		a, av = GetDeepValue(av, a)
		b, bv = GetDeepValue(bv, b)
		a = toNumberType(av)
		b = toNumberType(bv)
		return a.(float64) <= b.(float64), nil
	}
	return false, errors.New("find not support equal operator :" + operator)
}

func DoCalculationAction(operator Operator, a interface{}, b interface{}, av reflect.Value, bv reflect.Value) (interface{}, error) {
	if a == nil || b == nil {
		//equal nil
		return false, errors.New("add operator value can not be nil!")
	}
	//start equal
	a, av = GetDeepValue(av, a)
	b, bv = GetDeepValue(bv, b)
	switch operator {
	case Add:
		if av.Kind() == reflect.String {
			return a.(string) + b.(string), nil
		}
		a = toNumberType(av)
		b = toNumberType(bv)
		return a.(float64) + b.(float64), nil
	case Reduce:
		a = toNumberType(av)
		b = toNumberType(bv)
		return a.(float64) - b.(float64), nil
	case Ride:
		a = toNumberType(av)
		b = toNumberType(bv)
		return a.(float64) * b.(float64), nil
	case Divide:
		a = toNumberType(av)
		b = toNumberType(bv)
		if b.(float64) == 0 {
			return nil, errors.New("can not divide zero value!")
		}
		return a.(float64) / b.(float64), nil
	}
	return "", errors.New("find not support operator :" + operator)
}

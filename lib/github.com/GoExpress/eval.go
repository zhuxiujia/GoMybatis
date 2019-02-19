package GoExpress

import (
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
	if strings.Index(operator, ".") != -1 {
		if av.Kind() == reflect.Ptr {
			arg, av = GetDeepValue(av, arg)
		}
		var childs = strings.Split(operator, ".")
		var v reflect.Value
		for _, item := range childs {
			if av.Kind() == reflect.Struct {
				v = av.FieldByName(item)
				if v.Kind() == reflect.Ptr {
					v = v.Elem()
				}
			} else {
				panic(errors.New("only struct support take value a.b!"))
			}
		}
		if v.IsValid() && v.CanInterface() {
			return v.Interface(), nil
		} else {
			return nil, nil
		}
	} else {
		if av.Kind() == reflect.Ptr {
			arg, av = GetDeepValue(av, arg)
		}
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
	case Equal:
		//a kind == b kind
		return DoEqual(operator, a, b, av, bv)
	case UnEqual:
		//a kind == b kind
		var r, e = DoEqual(operator, a, b, av, bv)
		if e != nil {
			return nil, e
		}
		return !r, nil
	case Add:
		return DoAddReduceRideDivide(operator, a, b, av, bv)

	case Reduce:
		return DoAddReduceRideDivide(operator, a, b, av, bv)

		break
	case Ride:
		return DoAddReduceRideDivide(operator, a, b, av, bv)

		break
	case Divide:
		return DoAddReduceRideDivide(operator, a, b, av, bv)

		break
	}
	return nil, errors.New("find not support operator=" + operator)
}

func DoEqual(operator Operator, a interface{}, b interface{}, av reflect.Value, bv reflect.Value) (bool, error) {
	if a == nil || b == nil {
		//equal nil
		if a != nil || b != nil {
			return false, nil
		}
		if a == nil && b == nil {
			return true, nil
		}
	}

	//start equal
	a, av = GetDeepValue(av, a)
	b, bv = GetDeepValue(bv, b)

	if av.Kind() != bv.Kind() {
		return false, nil
	}
	if av.Kind() == reflect.Bool {
		return a.(bool) == b.(bool), nil
	}
	if av.Kind() == reflect.String {
		return a.(string) == b.(string), nil
	}
	a = toNumberType(av)
	b = toNumberType(bv)
	if isNumberType(av.Type()) {
		switch operator {
		case Equal:
			return a.(float64) == b.(float64), nil
		case UnEqual:
			return a.(float64) != b.(float64), nil
		case Less:
			return a.(float64) < b.(float64), nil
		case More:
			return a.(float64) > b.(float64), nil
		case MoreEqual:
			return a.(float64) >= b.(float64), nil
		case LessEqual:
			return a.(float64) <= b.(float64), nil
		}
	}
	return false, errors.New("find not support equal operator=" + operator)
}

func DoAddReduceRideDivide(operator Operator, a interface{}, b interface{}, av reflect.Value, bv reflect.Value) (interface{}, error) {
	if a == nil || b == nil {
		//equal nil
		return false, errors.New("add operator value can not be nil!")
	}
	//start equal
	a, av = GetDeepValue(av, a)
	b, bv = GetDeepValue(bv, b)

	if av.Kind() == reflect.String {
		return a.(string) + b.(string), nil
	}
	a = toNumberType(av)
	b = toNumberType(bv)
	if isNumberType(av.Type()) {
		switch operator {
		case Add:
			return a.(float64) + b.(float64), nil
		case Reduce:
			return a.(float64) - b.(float64), nil
		case Ride:
			return a.(float64) * b.(float64), nil
		case Divide:
			if b.(float64) == 0 {
				return nil, errors.New("can not divide zero value!")
			}
			return a.(float64) / b.(float64), nil
		}
	}
	return "", errors.New("find not support operator=" + operator)
}

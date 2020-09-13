package GoFastExpress

import (
	"errors"
	"fmt"
	"reflect"
)

//取值
func EvalTakes(argNode ArgNode, arg interface{}) (interface{}, error) {
	if arg == nil || argNode.values == nil {
		return nil, nil
	}
	if argNode.value == "" || argNode.valuesLen == 0 {
		return arg, nil
	}
	var av = reflect.ValueOf(arg)
	if av.Kind() == reflect.Map {
		var m = arg.(map[string]interface{})
		if argNode.valuesLen == 1 {
			return m[argNode.value], nil
		}
		return takeValue(argNode.value, av.MapIndex(reflect.ValueOf(argNode.values[0])), argNode.values[1:])
	} else {
		if argNode.valuesLen == 1 {
			return arg, nil
		}
		return takeValue(argNode.value, av, argNode.values[1:])
	}
}

func takeValue(key string, arg reflect.Value, feilds []string) (interface{}, error) {
	if arg.IsValid() == false {
		return nil, nil
	}
	for _, v := range feilds {
		argItem, e := getObjV(key, v, arg)
		if e != nil || argItem == nil {
			return nil, e
		}
		arg = *argItem
	}
	if !arg.IsValid() {
		return nil, nil
	}
	if arg.CanInterface() {
		var intf = arg.Interface()
		return intf, nil
	} else {
		return nil, nil
	}

}

func getObjV(key string, operator Operator, av reflect.Value) (*reflect.Value, error) {
	if av.Kind() == reflect.Ptr || av.Kind() == reflect.Interface {
		av = GetDeepPtr(av)
	}

	if av.Kind() == reflect.Map {
		var mapV = av.MapIndex(reflect.ValueOf(operator))
		return &mapV, nil
	}

	if av.Kind() != reflect.Struct {
		return nil, errors.New("[express] " + key + " get value  " + key + "  fail :" + av.String() + ",value key:" + operator)
	}
	av = av.FieldByName(operator)
	if av.Kind() == reflect.Ptr || av.Kind() == reflect.Interface {
		av = GetDeepPtr(av)
	}
	if av.IsValid() && av.CanInterface() {
		return &av, nil
	} else {
		return nil, nil
	}
}

func Eval(express string, operator Operator, a interface{}, b interface{}) (interface{}, error) {
	var av = reflect.ValueOf(a)
	var bv = reflect.ValueOf(b)

	switch operator {
	case And:
		if a == nil || b == nil {
			//equal nil
			return nil, errors.New("[express] " + express + " eval fail,value can not be nil")
		}
		a, av = GetDeepValue(av, a)
		b, bv = GetDeepValue(bv, b)
		var ab = a.(bool)
		var bb = b.(bool)
		return ab == true && bb == true, nil
	case Or:
		if a == nil || b == nil {
			//equal nil
			return nil, errors.New("[express] " + express + " eval fail,value can not be nil")
		}
		a, av = GetDeepValue(av, a)
		b, bv = GetDeepValue(bv, b)
		var ab = a.(bool)
		var bb = b.(bool)
		return ab == true || bb == true, nil
	case Equal, MoreEqual, More, Less, LessEqual:
		//a kind == b kind
		return DoEqualAction(express, operator, a, b, av, bv)
	case UnEqual:
		//a kind == b kind
		var r, e = DoEqualAction(express, operator, a, b, av, bv)
		if e != nil {
			return nil, e
		}
		return !r, nil
	case Add, Reduce, Ride, Divide:
		return DoCalculationAction(express, operator, a, b, av, bv)
	}
	return nil, errors.New("[express] " + express + " find not support operator :" + operator)
}

func DoEqualAction(express string, operator Operator, a interface{}, b interface{}, av reflect.Value, bv reflect.Value) (bool, error) {
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
			//empty map not equal nil
			if av.IsNil() && b == nil{
				return false,nil
			}
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
			return fmt.Sprintf(`%v`,a) == fmt.Sprintf(`%v`,b), nil
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
			return false, errors.New("[express] " + express + "can not parser '<' , arg have nil object!")
		}
		a, av = GetDeepValue(av, a)
		b, bv = GetDeepValue(bv, b)
		a = toNumberType(av)
		b = toNumberType(bv)
		return a.(float64) < b.(float64), nil
	case More:
		if a == nil || b == nil {
			return false, errors.New("[express] " + express + "can not parser '>' , arg have nil object!")
		}
		a, av = GetDeepValue(av, a)
		b, bv = GetDeepValue(bv, b)
		a = toNumberType(av)
		b = toNumberType(bv)
		return a.(float64) > b.(float64), nil
	case MoreEqual:
		if a == nil || b == nil {
			return false, errors.New("[express] " + express + "can not parser '>=' , arg have nil object!")
		}
		a, av = GetDeepValue(av, a)
		b, bv = GetDeepValue(bv, b)
		a = toNumberType(av)
		b = toNumberType(bv)
		return a.(float64) >= b.(float64), nil
	case LessEqual:
		if a == nil || b == nil {
			return false, errors.New("[express] " + express + "can not parser '<=' , arg have nil object!")
		}
		a, av = GetDeepValue(av, a)
		b, bv = GetDeepValue(bv, b)
		a = toNumberType(av)
		b = toNumberType(bv)
		return a.(float64) <= b.(float64), nil
	}
	return false, errors.New("[express] " + express + " find not support equal operator :" + operator)
}

func DoCalculationAction(express string, operator Operator, a interface{}, b interface{}, av reflect.Value, bv reflect.Value) (interface{}, error) {
	if a == nil || b == nil {
		//equal nil
		return false, errors.New("[express] " + express + " have not a action operator!")
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
			return nil, errors.New("[express] " + express + "can not divide zero value!")
		}
		return a.(float64) / b.(float64), nil
	}
	return "", errors.New("[express] " + express + "find not support operator :" + operator)
}

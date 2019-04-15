package utils

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func GetFloatValuePtr(arg *float64) string {
	if arg != nil {
		//strconv.FormatFloat(*arg, )
	}
	return ""
}

func GetIntValuePtr(arg *int) string {
	if arg != nil {
		return strconv.Itoa(*arg)
	}
	return ""
}

func GetStringValuePtr(arg *string) string {
	if arg != nil {
		return *arg
	}
	return ""
}

func GetTimeValuePtr(arg *time.Time) string {
	if arg != nil {
		return (*arg).Format(time.RFC3339)
	}
	return ""
}

func GetValue(arg interface{}, types reflect.Type) string {
	if arg != nil {
		if types != nil {
			if types.Kind() == reflect.Ptr {
				return caseTypePtr(arg, types)
			} else {
				return caseType(arg, types)
			}
		}
		var v = reflect.ValueOf(arg)
		if v.Type().Kind() == reflect.Ptr {
			v = v.Elem()
		}
		switch v.Kind() {
		case reflect.String:
			return v.String()
		case reflect.Int:
			return strconv.FormatInt(v.Int(), 36)
		case reflect.Int8:
			return strconv.FormatInt(v.Int(), 8)
		case reflect.Int16:
			return strconv.FormatInt(v.Int(), 16)
		case reflect.Int32:
			return strconv.FormatInt(v.Int(), 32)
		case reflect.Int64:
			return strconv.FormatInt(v.Int(), 36)
		case reflect.Float32:
			return strconv.FormatFloat(v.Float(), 'f', -1, 32)
		case reflect.Float64:
			return strconv.FormatFloat(v.Float(), 'f', -1, 64)
		case reflect.Bool:
			return strconv.FormatBool(v.Bool())
		case reflect.Uint:
			return strconv.FormatUint(v.Uint(), 2)
		case reflect.Uint8:
			return strconv.FormatUint(v.Uint(), 8)
		case reflect.Uint16:
			return strconv.FormatUint(v.Uint(), 16)
		case reflect.Uint32:
			return strconv.FormatUint(v.Uint(), 32)
		case reflect.Uint64:
			return strconv.FormatUint(v.Uint(), 36)
		case reflect.Struct:
			if strings.Index(v.String(), "time.Time") != -1 {
				return v.Interface().(time.Time).Format(`2006-01-02 15:04:05`)
			}
		default:
			return v.String()
		}
	}
	return ""
}

func caseType(arg interface{}, types reflect.Type) string {
	switch types.Kind() {
	case reflect.String:
		return arg.(string)
	case reflect.Int:
		return strconv.FormatInt(int64(arg.(int)), 10)
	case reflect.Int8:
		return strconv.FormatInt(int64(arg.(int8)), 10)
	case reflect.Int16:
		return strconv.FormatInt(int64(arg.(int16)), 10)
	case reflect.Int32:
		return strconv.FormatInt(int64(arg.(int32)), 10)
	case reflect.Int64:
		return strconv.FormatInt(int64(arg.(int64)), 10)
	case reflect.Float32:
		return strconv.FormatFloat(float64(arg.(float32)), 'f', -1, 64)
	case reflect.Float64:
		return strconv.FormatFloat(float64(arg.(float64)), 'f', -1, 64)
	case reflect.Bool:
		return strconv.FormatBool(arg.(bool))
	case reflect.Uint:
		return strconv.FormatUint(uint64(arg.(uint)), 10)
	case reflect.Uint8:
		return strconv.FormatUint(uint64(arg.(uint8)), 10)
	case reflect.Uint16:
		return strconv.FormatUint(uint64(arg.(uint16)), 10)
	case reflect.Uint32:
		return strconv.FormatUint(uint64(arg.(uint32)), 10)
	case reflect.Uint64:
		return strconv.FormatUint(uint64(arg.(uint64)), 10)
	case reflect.Struct:
		if types.String() == "time.Time" {
			return arg.(time.Time).Format(`2006-01-02 15:04:05`)
		} else {
			return fmt.Sprint(arg)
		}
	default:
		return fmt.Sprint(arg)
	}
	return fmt.Sprint(arg)
}

func caseTypePtr(arg interface{}, types reflect.Type) string {
	var childType = types.Elem()
	switch childType.Kind() {
	case reflect.String:
		arg := arg.(*string)
		if arg == nil {
			return ""
		}
		return *arg
	case reflect.Int:
		arg := arg.(*int)
		if arg == nil {
			return ""
		}
		return strconv.FormatInt(int64(*arg), 8)
	case reflect.Int8:
		arg := arg.(*int8)
		if arg == nil {
			return ""
		}
		return strconv.FormatInt(int64(*arg), 8)
	case reflect.Int16:
		arg := arg.(*int16)
		if arg == nil {
			return ""
		}
		return strconv.FormatInt(int64(*arg), 16)
	case reflect.Int32:
		arg := arg.(*int32)
		if arg == nil {
			return ""
		}
		return strconv.FormatInt(int64(*arg), 32)
	case reflect.Int64:
		arg := arg.(*int64)
		if arg == nil {
			return ""
		}
		return strconv.FormatInt(int64(*arg), 64)
	case reflect.Float32:
		arg := arg.(*float32)
		if arg == nil {
			return ""
		}
		return strconv.FormatFloat(float64(*arg), 'f', -1, 32)
	case reflect.Float64:
		arg := arg.(*float64)
		if arg == nil {
			return ""
		}
		return strconv.FormatFloat(float64(*arg), 'f', -1, 64)
	case reflect.Bool:
		arg := arg.(*bool)
		if arg == nil {
			return ""
		}
		return strconv.FormatBool(*arg)
	case reflect.Uint:
		arg := arg.(*uint)
		if arg == nil {
			return ""
		}
		return strconv.FormatUint(uint64(*arg), 2)
	case reflect.Uint8:
		arg := arg.(*uint8)
		if arg == nil {
			return ""
		}
		return strconv.FormatUint(uint64(*arg), 8)
	case reflect.Uint16:
		arg := arg.(*uint16)
		if arg == nil {
			return ""
		}
		return strconv.FormatUint(uint64(*arg), 16)
	case reflect.Uint32:
		arg := arg.(*uint32)
		if arg == nil {
			return ""
		}
		return strconv.FormatUint(uint64(*arg), 32)
	case reflect.Uint64:
		arg := arg.(*uint64)
		if arg == nil {
			return ""
		}
		return strconv.FormatUint(uint64(*arg), 64)
	case reflect.Struct:
		if types.String() == "*time.Time" {
			arg := arg.(*time.Time)
			if arg == nil {
				return ""
			}
			return (*arg).Format(`2006-01-02 15:04:05`)
		} else {
			return fmt.Sprint(arg)
		}
	default:
		return fmt.Sprint(arg)
	}
	return fmt.Sprint(arg)
}

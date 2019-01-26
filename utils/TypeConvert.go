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
		return strconv.FormatInt(int64(arg.(int)), 8)
	case reflect.Int8:
		return strconv.FormatInt(int64(arg.(int8)), 8)
	case reflect.Int16:
		return strconv.FormatInt(int64(arg.(int16)), 16)
	case reflect.Int32:
		return strconv.FormatInt(int64(arg.(int32)), 32)
	case reflect.Int64:
		return strconv.FormatInt(int64(arg.(int64)), 64)
	case reflect.Float32:
		return strconv.FormatFloat(float64(arg.(float32)), 'f', -1, 32)
	case reflect.Float64:
		return strconv.FormatFloat(float64(arg.(float64)), 'f', -1, 64)
	case reflect.Bool:
		return strconv.FormatBool(arg.(bool))
	case reflect.Uint:
		return strconv.FormatUint(uint64(arg.(uint)), 2)
	case reflect.Uint8:
		return strconv.FormatUint(uint64(arg.(uint8)), 8)
	case reflect.Uint16:
		return strconv.FormatUint(uint64(arg.(uint16)), 16)
	case reflect.Uint32:
		return strconv.FormatUint(uint64(arg.(uint32)), 32)
	case reflect.Uint64:
		return strconv.FormatUint(uint64(arg.(uint64)), 64)
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
		return *arg.(*string)
	case reflect.Int:
		return strconv.FormatInt(int64(*arg.(*int)), 8)
	case reflect.Int8:
		return strconv.FormatInt(int64(*arg.(*int8)), 8)
	case reflect.Int16:
		return strconv.FormatInt(int64(*arg.(*int16)), 16)
	case reflect.Int32:
		return strconv.FormatInt(int64(*arg.(*int32)), 32)
	case reflect.Int64:
		return strconv.FormatInt(int64(*arg.(*int64)), 64)
	case reflect.Float32:
		return strconv.FormatFloat(float64(*arg.(*float32)), 'f', -1, 32)
	case reflect.Float64:
		return strconv.FormatFloat(float64(*arg.(*float64)), 'f', -1, 64)
	case reflect.Bool:
		return strconv.FormatBool(*arg.(*bool))
	case reflect.Uint:
		return strconv.FormatUint(uint64(*arg.(*uint)), 2)
	case reflect.Uint8:
		return strconv.FormatUint(uint64(*arg.(*uint8)), 8)
	case reflect.Uint16:
		return strconv.FormatUint(uint64(*arg.(*uint16)), 16)
	case reflect.Uint32:
		return strconv.FormatUint(uint64(*arg.(*uint32)), 32)
	case reflect.Uint64:
		return strconv.FormatUint(uint64(*arg.(*uint64)), 64)
	case reflect.Struct:
		if types.String() == "*time.Time" {
			return (*arg.(*time.Time)).Format(`2006-01-02 15:04:05`)
		} else {
			return fmt.Sprint(arg)
		}
	default:
		return fmt.Sprint(arg)
	}
	return fmt.Sprint(arg)
}

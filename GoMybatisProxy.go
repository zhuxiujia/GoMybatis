package GoMybatis

import (
	"reflect"
	"strings"
)

// UseService 可写入每个函数代理方法
func UseMapper(mapper interface{}, proxyFunc func(method string, args []reflect.Value, params []string) error) {
	v := reflect.ValueOf(mapper)
	if v.Kind() != reflect.Ptr {
		panic("UseMapper: UseMapper arg must be a pointer")
	}
	buildMapper(v, proxyFunc)
}

// UseService 可写入每个函数代理方法
func UseMapperValue(mapperValue reflect.Value, proxyFunc func(method string, args []reflect.Value, params []string) error) {
	buildMapper(mapperValue, proxyFunc)
}

func buildMapper(v reflect.Value, proxyFunc func(method string, args []reflect.Value, params []string) error) {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()
	et := t
	if et.Kind() == reflect.Ptr {
		et = et.Elem()
	}
	ptr := reflect.New(et)
	obj := ptr.Elem()
	count := obj.NumField()
	for i := 0; i < count; i++ {
		f := obj.Field(i)
		ft := f.Type()
		sf := et.Field(i)
		if ft.Kind() == reflect.Ptr {
			ft = ft.Elem()
		}
		if f.CanSet() {
			switch ft.Kind() {
			case reflect.Struct:
			case reflect.Func:
				buildRemoteMethod(f, ft, sf, proxyFunc)
			}
		}
	}
	if t.Kind() == reflect.Ptr {
		v.Set(ptr)
	} else {
		v.Set(obj)
	}
}

func buildRemoteMethod(f reflect.Value, ft reflect.Type, sf reflect.StructField, proxyFunc func(method string, args []reflect.Value, params []string) error) {
	var params []string
	var mapperParams = sf.Tag.Get(`mapperParams`)
	if mapperParams != `` {
		params = strings.Split(mapperParams, `,`)
	}
	var fn = func(args []reflect.Value) (results []reflect.Value) {
		err := proxyFunc(sf.Name, args, params)
		results = append(results, reflect.ValueOf(&err).Elem())
		return
	}
	if f.Kind() == reflect.Ptr {
		fp := reflect.New(ft)
		fp.Elem().Set(reflect.MakeFunc(ft, fn))
		f.Set(fp)
	} else {
		f.Set(reflect.MakeFunc(ft, fn))
	}
}

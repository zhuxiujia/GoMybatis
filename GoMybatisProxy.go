package GoMybatis

import (
	"reflect"
	"strings"
)

type TagArg struct {
	Name  string
	Index int
}

// UseService 可写入每个函数代理方法
func UseMapper(mapper interface{}, proxyFunc func(method string, args []reflect.Value, tagArgs []TagArg) []reflect.Value) {
	v := reflect.ValueOf(mapper)
	if v.Kind() != reflect.Ptr {
		panic("UseMapper: UseMapper arg must be a pointer")
	}
	buildMapper(v, proxyFunc)
}

// UseService 可写入每个函数代理方法
func UseMapperValue(mapperValue reflect.Value, proxyFunc func(method string, args []reflect.Value, tagArgs []TagArg) []reflect.Value) {
	buildMapper(mapperValue, proxyFunc)
}

func buildMapper(v reflect.Value, proxyFunc func(method string, args []reflect.Value, tagArgs []TagArg) []reflect.Value) {
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

func buildRemoteMethod(f reflect.Value, ft reflect.Type, sf reflect.StructField, proxyFunc func(method string, args []reflect.Value, tagArgs []TagArg) []reflect.Value) {
	var params []string
	var mapperParams = sf.Tag.Get(`mapperParams`)
	if mapperParams != `` {
		params = strings.Split(mapperParams, `,`)
	}
	if len(params) > ft.NumIn() {
		panic(`[GoMybatisProxy] method fail! the tag "mapperParams" length can not > arg length! filed=` + ft.String())
	}
	var tagArgs = make([]TagArg, 0)
	if len(params) != 0 {
		for index, v := range params {
			var tagArg = TagArg{
				Index: index,
				Name:  v,
			}
			tagArgs = append(tagArgs, tagArg)
		}
	}
	var fn = func(args []reflect.Value) (results []reflect.Value) {
		proxyResults := proxyFunc(sf.Name, args, tagArgs)
		for _, returnV := range proxyResults {
			results = append(results, returnV)
		}
		return results
	}
	if f.Kind() == reflect.Ptr {
		fp := reflect.New(ft)
		fp.Elem().Set(reflect.MakeFunc(ft, fn))
		f.Set(fp)
	} else {
		f.Set(reflect.MakeFunc(ft, fn))
	}
	params = nil
}

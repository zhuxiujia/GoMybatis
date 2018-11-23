package GoMybatis

import (
	"errors"
	"reflect"
	"strings"
)

//如果使用UseProxyMapperByEngine，则内建默认的SessionFactory
var DefaultSessionFactory *SessionFactory

var DefaultSqlResultDecoder SqlResultDecoder

var DefaultSqlBuilder SqlBuilder

func UseProxyMapperByEngine(bean interface{}, xml []byte, sqlEngine *SessionEngine) {
	v := reflect.ValueOf(bean)
	if v.Kind() != reflect.Ptr {
		panic("UseMapper: UseMapper arg must be a pointer")
	}
	if DefaultSessionFactory == nil {
		var factory = SessionFactory{}.New(sqlEngine)
		DefaultSessionFactory = &factory
	}
	if DefaultSqlResultDecoder == nil {
		DefaultSqlResultDecoder = GoMybatisSqlResultDecoder{}
	}
	if DefaultSqlBuilder == nil {
		DefaultSqlBuilder = GoMybatisSqlBuilder{}
	}
	UseProxyMapper(v, xml, DefaultSessionFactory, DefaultSqlResultDecoder, DefaultSqlBuilder)
}

func UseProxyMapperByFactory(bean interface{}, xml []byte, sessionFactory *SessionFactory, sqlResultDecoder SqlResultDecoder, sqlBuilder SqlBuilder) {
	v := reflect.ValueOf(bean)
	if v.Kind() != reflect.Ptr {
		panic("UseMapper: UseMapper arg must be a pointer")
	}
	UseProxyMapperFromValue(v, xml, sessionFactory, sqlResultDecoder, sqlBuilder)
}

func UseProxyMapperFromValue(bean reflect.Value, xml []byte, sessionFactory *SessionFactory, sqlResultDecoder SqlResultDecoder, sqlBuilder SqlBuilder) {
	UseProxyMapper(bean, xml, sessionFactory, sqlResultDecoder, sqlBuilder)
}

//例如
//type ExampleActivityMapperImpl struct {
//	SelectAll         func(result *[]Activity) error
//	SelectByCondition func(name string, startTime time.Time, endTime time.Time, page int, size int, result *[]Activity) error `mapperParams:"name,startTime,endTime,page,size"`
//	UpdateById        func(session *GoMybatis.Session, arg Activity, result *int64) error                                     //只要参数中包含有*GoMybatis.Session的类型，框架默认使用传入的session对象，用于自定义事务
//	Insert            func(arg Activity, result *int64) error
//	CountByCondition  func(name string, startTime time.Time, endTime time.Time, result *int) error `mapperParams:"name,startTime,endTime"`
//}
//bean 工厂，根据xml配置创建函数,并且动态代理到你定义的struct func里
//bean 参数必须为reflect.Value
//func的基本类型的参数（例如string,int,time.Time,int64,float....）个数无限制(并且需要用Tag指定参数名逗号隔开,例如`mapperParams:"id,phone"`)，最后一个参数必须为返回数据类型的指针(例如result *model.User)，返回值为error
//func的结构体参数无需指定mapperParams的tag，框架会自动扫描它的属性，封装为map处理掉
//使用UseProxyMapper函数设置代理后即可正常使用。
func UseProxyMapper(bean reflect.Value, xml []byte, sessionFactory *SessionFactory, decoder SqlResultDecoder, sqlBuilder SqlBuilder) {
	var mapperTree = LoadMapperXml(xml)
	var proxyFunc = func(method string, args []reflect.Value, tagArgs []TagArg) error {
		var lastArgsIndex = len(args) - 1
		var argsLen = len(args)
		var lastArgValue *reflect.Value = nil
		if argsLen != 0 && args[lastArgsIndex].Kind() == reflect.Ptr {
			lastArgValue = &args[lastArgsIndex]
			if lastArgValue.Kind() != reflect.Ptr {
				//最后一个参数必须为指针，或者不传任何参数
				return errors.New(`[method params last param must be pointer!],method =` + method)
			}
		}
		var findMethod = false
		for _, mapperXml := range mapperTree {
			//exec sql,return data
			if strings.EqualFold(mapperXml.Id, method) {
				findMethod = true
				buildMethodBody(sessionFactory, tagArgs, args, mapperXml, lastArgValue, decoder, sqlBuilder)
				//匹配完成退出
				break
			}
		}
		if findMethod == false {
			return errors.New(`[GoMybatis] not method find at xml file,method =` + method)
		}
		return nil
	}
	UseMapperValue(bean, proxyFunc)
}

func buildMethodBody(sessionFactory *SessionFactory, tagParamMap []TagArg, args []reflect.Value, mapperXml MapperXml, lastArgValue *reflect.Value, decoder SqlResultDecoder, sqlBuilder SqlBuilder) error {
	//build sql string
	var session *Session
	var sql string
	var err error
	session, sql, err = buildSql(tagParamMap, args, mapperXml, sqlBuilder)
	if err != nil {
		return err
	}
	//session
	var haveArgSession = false
	if session == nil {
		session = sessionFactory.NewSession()
	} else {
		haveArgSession = true
	}
	if haveArgSession == false {
		//not arg session,just close!
		defer closeSession(sessionFactory, session)
	}
	var haveLastReturnValue = lastArgValue != nil && (*lastArgValue).IsNil() == false
	//do CRUD
	if mapperXml.Tag == Select && haveLastReturnValue {
		//is select and have return value
		results, err := (*session).Query(sql)
		if err != nil {
			return err
		}
		err = decoder.Decode(results, lastArgValue.Interface())
		if err != nil {
			return err
		}
	} else {
		var res, err = (*session).Exec(sql)
		if err != nil {
			return err
		}
		if haveLastReturnValue {
			lastArgValue.Elem().SetInt(res.RowsAffected)
		}
	}
	return nil
}

func closeSession(factory *SessionFactory, session *Session) {
	factory.CloseSession((*session).Id())
	(*session).Close()
}

func buildSql(tagArgs []TagArg, args []reflect.Value, mapperXml MapperXml, sqlBuilder SqlBuilder) (*Session, string, error) {
	var session *Session
	var paramMap = make(map[string]interface{})
	var tagArgsLen = len(tagArgs)
	for argIndex, arg := range args {
		var argInterface = arg.Interface()
		if arg.Kind() == reflect.Ptr {
			if arg.Type().String() == GoMybatis_Session {
				//指针，则退出
				session = argInterface.(*Session)
			}
			continue
		}
		if arg.Kind() == reflect.Struct && arg.Type().String() != GoMybatis_Time {
			paramMap = scanStructArgFields(argInterface, nil)
		} else if tagArgsLen > 0 && argIndex < tagArgsLen && tagArgs[argIndex].Name != "" && argInterface != nil {
			paramMap[tagArgs[argIndex].Name] = argInterface
		} else {
			paramMap[DefaultOneArg] = argInterface
		}
	}
	result, err := sqlBuilder.BuildSqlFromMap(paramMap, mapperXml)
	return session, result, err
}

//scan params
func scanStructArgFields(arg interface{}, typeConvert func(arg interface{}) interface{}) map[string]interface{} {
	parameters := make(map[string]interface{})
	v := reflect.ValueOf(arg)
	t := reflect.TypeOf(arg)
	if t.Kind() != reflect.Struct {
		panic(`[GoMybatis] the scanParamterBean() arg is not a struct type!,type =` + t.String())
	}
	for i := 0; i < t.NumField(); i++ {
		var typeValue = t.Field(i)
		var obj = v.Field(i).Interface()
		if typeConvert != nil {
			obj = typeConvert(obj)
		}
		var jsonKey = typeValue.Tag.Get(`json`)
		if jsonKey != "" {
			parameters[jsonKey] = obj
		} else {
			parameters[typeValue.Name] = obj
		}
	}
	return parameters
}

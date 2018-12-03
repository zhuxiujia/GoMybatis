package GoMybatis

import (
	"errors"
	"reflect"
	"strings"
)

//如果使用UseProxyMapperByEngine，则内建默认的SessionFactory
var DefaultSessionFactory *SessionFactory

func UseProxyMapperByEngine(bean interface{}, xml []byte, sqlEngine *SessionEngine) {
	v := reflect.ValueOf(bean)
	if v.Kind() != reflect.Ptr {
		panic("UseMapper: UseMapper arg must be a pointer")
	}
	var factory = SessionFactory{}.New(sqlEngine)
	if DefaultSessionFactory == nil {
		DefaultSessionFactory = &factory
	}
	UseProxyMapper(v, xml, DefaultSessionFactory, GoMybatisSqlResultDecoder{}, GoMybatisSqlBuilder{}.New(GoMybatisExpressionTypeConvert{}, GoMybatisSqlArgTypeConvert{}))
}

func UseProxyMapperFromBean(bean interface{}, xml []byte, sessionFactory *SessionFactory, sqlResultDecoder SqlResultDecoder, sqlBuilder SqlBuilder) {
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
	//make a map[method]xml
	var methodXmlMap = makeMethodXmlMap(bean, mapperTree)
	var resultMaps = makeResultMaps(mapperTree)
	var proxyFunc = func(method string, args []reflect.Value, tagArgs []TagArg) error {
		var lastArgsIndex = len(args) - 1
		var argsLen = len(args)
		var lastArgValue *reflect.Value = nil
		if argsLen != 0 && args[lastArgsIndex].Kind() == reflect.Ptr {
			lastArgValue = &args[lastArgsIndex]
			if lastArgValue.Kind() != reflect.Ptr {
				//最后一个参数必须为指针，或者不传任何参数
				return errors.New(`[GoMybatis] method params last param must be pointer!,method =` + method)
			}
		}
		var mapperXml = methodXmlMap[method]
		var resultMap map[string]*ResultProperty
		var resultMapId = mapperXml.Propertys["resultMap"]
		if resultMapId != "" {
			resultMap = resultMaps[resultMapId]
		}
		return exeMethodByXml(sessionFactory, tagArgs, args, mapperXml, resultMap, lastArgValue, decoder, sqlBuilder)
	}
	UseMapperValue(bean, proxyFunc)
}

//map[id]map[cloum]Property
func makeResultMaps(xmls map[string]*MapperXml) map[string]map[string]*ResultProperty {
	var resultMaps = make(map[string]map[string]*ResultProperty)
	for _, xmlItem := range xmls {
		if xmlItem.Tag == "resultMap" {
			var resultPropertyMap = make(map[string]*ResultProperty)
			for _, elementItem := range xmlItem.ElementItems {
				var property = ResultProperty{
					XMLName:  elementItem.ElementType,
					Column:   elementItem.Propertys["column"],
					Property: elementItem.Propertys["property"],
					GoType:   elementItem.Propertys["goType"],
				}
				resultPropertyMap[property.Column] = &property
			}
			resultMaps[xmlItem.Id] = resultPropertyMap
		}
	}
	return resultMaps
}
//return a map map[`method`]*MapperXml
func makeMethodXmlMap(bean reflect.Value, mapperTree map[string]*MapperXml) map[string]*MapperXml {
	if bean.Kind() == reflect.Ptr {
		bean = bean.Elem()
	}

	var methodXmlMap = make(map[string]*MapperXml)
	var totalField = bean.Type().NumField()
	for i := 0; i < totalField; i++ {
		var fieldItem = bean.Type().Field(i)
		if fieldItem.Type.Kind() == reflect.Func {
			//field must be func
			methodFieldCheck(fieldItem)
			var mapperXml = findMapperXml(mapperTree, fieldItem.Name)
			if mapperXml != nil {
				methodXmlMap[fieldItem.Name] = mapperXml
			} else {
				panic("[GoMybatis] can not find method "+bean.Type().String() +"."+ fieldItem.Name + "() in xml !")
			}
		}
	}
	return methodXmlMap
}

func methodFieldCheck(methodType reflect.StructField) {
	if methodType.Type.NumOut() != 1 {
		panic("[GoMybatis] method field must be return one 'error' type!")
	}
}

func findMapperXml(mapperTree map[string]*MapperXml, methodName string) *MapperXml {
	for _, mapperXml := range mapperTree {
		//exec sql,return data
		if strings.EqualFold(mapperXml.Id, methodName) {
			return mapperXml
		}
	}
	return nil
}

func exeMethodByXml(sessionFactory *SessionFactory, tagParamMap []TagArg, args []reflect.Value, mapperXml *MapperXml, resultMap map[string]*ResultProperty, lastArgValue *reflect.Value, decoder SqlResultDecoder, sqlBuilder SqlBuilder) error {
	//build sql string
	var session *Session
	var sql string
	var err error
	session, sql, err = buildSql(tagParamMap, args, mapperXml, sqlBuilder)
	if err != nil {
		return err
	}
	//session
	if session == nil {
		session = sessionFactory.NewSession()
		//not arg session,just close!
		defer closeSession(sessionFactory, session)
	}
	var haveLastReturnValue = lastArgValue != nil && (*lastArgValue).IsNil() == false
	//do CRUD
	if mapperXml.Tag == Element_Select && haveLastReturnValue {
		//is select and have return value
		results, err := (*session).Query(sql)
		if err != nil {
			return err
		}
		err = decoder.Decode(resultMap, results, lastArgValue.Interface())
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

func buildSql(tagArgs []TagArg, args []reflect.Value, mapperXml *MapperXml, sqlBuilder SqlBuilder) (*Session, string, error) {
	var session *Session
	var paramMap = make(map[string]SqlArg)
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
			paramMap[tagArgs[argIndex].Name] = SqlArg{
				Value: argInterface,
				Type:  arg.Type(),
			}
		} else {
			paramMap[DefaultOneArg] = SqlArg{
				Value: argInterface,
				Type:  arg.Type(),
			}
		}
	}
	result, err := sqlBuilder.BuildSql(paramMap, mapperXml)
	return session, result, err
}

//scan params
func scanStructArgFields(arg interface{}, typeConvert func(arg interface{}) interface{}) map[string]SqlArg {
	parameters := make(map[string]SqlArg)
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
			parameters[jsonKey] = SqlArg{
				Type:  v.Field(i).Type(),
				Value: obj,
			}
		} else {
			parameters[typeValue.Name] = SqlArg{
				Type:  v.Field(i).Type(),
				Value: obj,
			}
		}
	}
	return parameters
}

package GoMybatis

import (
	"reflect"
	"strings"
)

//如果使用WriteMapperByEngine()，则内建默认的SessionFactory
var DefaultSessionFactory *SessionFactory

//根据sessionEngine写入到mapperPtr
func WriteMapperByEngine(value reflect.Value, xml []byte, sessionEngine *SessionEngine, enableLog bool) {
	if value.Kind() != reflect.Ptr {
		panic("UseMapper: UseMapper arg must be a pointer")
	}
	var factory = SessionFactory{}.New(sessionEngine)
	if DefaultSessionFactory == nil {
		DefaultSessionFactory = &factory
	}
	WriteMapper(value, xml, DefaultSessionFactory, GoMybatisSqlResultDecoder{}, GoMybatisSqlBuilder{}.New(GoMybatisExpressionTypeConvert{}, GoMybatisSqlArgTypeConvert{}), enableLog)
}

//根据sessionEngine写入到mapperPtr
func WriteMapperPtrByEngine(ptr interface{}, xml []byte, sessionEngine *SessionEngine, enableLog bool) {
	v := reflect.ValueOf(ptr)
	if v.Kind() != reflect.Ptr {
		panic("UseMapper: UseMapper arg must be a pointer")
	}
	WriteMapperByEngine(v, xml, sessionEngine, enableLog)
}

//写入方法内容，例如
//type ExampleActivityMapperImpl struct {
//	SelectAll         func(result *[]Activity) error
//	SelectByCondition func(name string, startTime time.Time, endTime time.Time, page int, size int, result *[]Activity) error `mapperParams:"name,startTime,endTime,page,size"`
//	UpdateById        func(session *GoMybatis.Session, arg Activity, result *int64) error                                     //只要参数中包含有*GoMybatis.Session的类型，框架默认使用传入的session对象，用于自定义事务
//	Insert            func(arg Activity, result *int64) error
//	CountByCondition  func(name string, startTime time.Time, endTime time.Time, result *int) error `mapperParams:"name,startTime,endTime"`
//}
//func的基本类型的参数（例如string,int,time.Time,int64,float....）个数无限制(并且需要用Tag指定参数名逗号隔开,例如`mapperParams:"id,phone"`)，返回值必须有error
//func的结构体参数无需指定mapperParams的tag，框架会自动扫描它的属性，封装为map处理掉
//使用WriteMapper函数设置代理后即可正常使用。
func WriteMapper(bean reflect.Value, xml []byte, sessionFactory *SessionFactory, decoder SqlResultDecoder, sqlBuilder SqlBuilder, enableLog bool) {
	var mapperTree = LoadMapperXml(xml)
	//make a map[method]xml
	var methodXmlMap = makeMethodXmlMap(bean, mapperTree)
	var resultMaps = makeResultMaps(mapperTree)
	var returnTypeMap = makeReturnTypeMap(bean)
	var proxyFunc = func(method string, args []reflect.Value, tagArgs []TagArg) []reflect.Value {
		var returnValue *reflect.Value = nil
		var returnType = returnTypeMap[method]

		if returnType == nil {
			panic("[GoMybatis] struct have no return values!")
		}

		if returnType.ReturnOutType != nil {
			var returnV = reflect.New(*returnType.ReturnOutType)
			switch (*returnType.ReturnOutType).Kind() {
			case reflect.Map:
				returnV.Elem().Set(reflect.MakeMap(*returnType.ReturnOutType))
			case reflect.Slice:
				returnV.Elem().Set(reflect.MakeSlice(*returnType.ReturnOutType, 0, 0))
			}
			returnValue = &returnV
		}

		var mapperXml = methodXmlMap[method]
		var resultMap map[string]*ResultProperty
		var resultMapId = mapperXml.Propertys[Element_ResultMap]
		if resultMapId != "" {
			resultMap = resultMaps[resultMapId]
		}

		var e = exeMethodByXml(sessionFactory, tagArgs, args, mapperXml, resultMap, returnValue, decoder, sqlBuilder, enableLog)

		var returnValues = make([]reflect.Value, returnType.NumOut)
		for index, _ := range returnValues {
			if index == returnType.ReturnIndex {
				if returnValue != nil {
					returnValues[index] = (*returnValue).Elem()
				}
			} else {
				if e != nil {
					returnValues[index] = reflect.New(*returnType.ErrorType)
					returnValues[index].Elem().Set(reflect.ValueOf(e))
					returnValues[index] = returnValues[index].Elem()
				} else {
					returnValues[index] = reflect.Zero(*returnType.ErrorType)
				}
			}
		}
		return returnValues
	}
	UseMapperValue(bean, proxyFunc)
}

func makeReturnTypeMap(value reflect.Value) (returnMap map[string]*ReturnType) {
	returnMap = make(map[string]*ReturnType)
	var proxyType = value.Elem().Type()
	for i := 0; i < proxyType.NumField(); i++ {
		var funcType = proxyType.Field(i).Type
		var key = proxyType.Field(i).Name
		var numOut = funcType.NumOut()
		if numOut > 2 || numOut == 0 {
			panic("[GoMybatis] func num out must = 1 or = 2,and have return a error!")
		}
		for f := 0; f < numOut; f++ {
			var outType = funcType.Out(f)
			if outType.Kind() == reflect.Ptr {
				panic("[GoMybatis] func return can not be a ptr!")
			}
			var returnType = returnMap[key]
			if returnType == nil {
				returnMap[key] = &ReturnType{
					ReturnIndex: -1,
					NumOut:      numOut,
				}
			}
			if outType.String() != "error" {
				returnMap[key].ReturnIndex = f
				returnMap[key].ReturnOutType = &outType
			} else {
				returnMap[key].ErrorType = &outType
			}
		}
	}
	return returnMap
}

//map[id]map[cloum]Property
func makeResultMaps(xmls map[string]*MapperXml) map[string]map[string]*ResultProperty {
	var resultMaps = make(map[string]map[string]*ResultProperty)
	for _, xmlItem := range xmls {
		if xmlItem.Tag == Element_ResultMap {
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
				panic("[GoMybatis] can not find method " + bean.Type().String() + "." + fieldItem.Name + "() in xml !")
			}
		}
	}
	return methodXmlMap
}

func methodFieldCheck(methodType reflect.StructField) {
	if methodType.Type.NumOut() < 1 {
		panic("[GoMybatis] method " + methodType.Name + "() must be return a 'error' type!")
	}
	var errorTypeNum = 0
	for i := 0; i < methodType.Type.NumOut(); i++ {
		var outType = methodType.Type.Out(i)
		if outType.Kind() == reflect.Interface && outType.String() == "error" {
			errorTypeNum++
		}
	}
	if errorTypeNum != 1 {
		panic("[GoMybatis] method " + methodType.Name + "() must be return a 'error' type!")
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

func exeMethodByXml(sessionFactory *SessionFactory, tagParamMap []TagArg, args []reflect.Value, mapperXml *MapperXml, resultMap map[string]*ResultProperty, returnValue *reflect.Value, decoder SqlResultDecoder, sqlBuilder SqlBuilder, enableLog bool) error {
	//build sql string
	var session Session
	var sql string
	var err error
	session, sql, err = buildSql(tagParamMap, args, mapperXml, sqlBuilder, enableLog)
	if err != nil {
		return err
	}
	if sessionFactory == nil && session == nil {
		panic("[GoMybatis] exe sql need a SessionFactory or Session!")
	}
	//session
	if session == nil {
		session = sessionFactory.NewSession(SessionType_Default, nil)
		//not arg session,just close!
		defer closeSession(sessionFactory, session)
	}
	var haveLastReturnValue = returnValue != nil && (*returnValue).IsNil() == false
	//do CRUD
	if mapperXml.Tag == Element_Select && haveLastReturnValue {
		//is select and have return value
		results, err := session.Query(sql)
		if err != nil {
			return err
		}
		err = decoder.Decode(resultMap, results, returnValue.Interface())
		if err != nil {
			return err
		}
	} else {
		var res, err = session.Exec(sql)
		if err != nil {
			return err
		}
		if haveLastReturnValue {
			returnValue.Elem().SetInt(res.RowsAffected)
		}
	}
	return nil
}

func closeSession(factory *SessionFactory, session Session) {
	if session == nil {
		return
	}
	factory.Close(session.Id())
	session.Close()
}

func buildSql(tagArgs []TagArg, args []reflect.Value, mapperXml *MapperXml, sqlBuilder SqlBuilder, enableLog bool) (Session, string, error) {
	var session Session
	var paramMap = make(map[string]SqlArg)
	var tagArgsLen = len(tagArgs)
	for argIndex, arg := range args {
		var argInterface = arg.Interface()
		if arg.Kind() == reflect.Ptr && arg.IsNil() == false && argInterface != nil && arg.Type().String() == GoMybatis_Session_Ptr {
			session = *(argInterface.(*Session))
			continue
		} else if argInterface != nil && arg.Kind() == reflect.Interface && arg.Type().String() == GoMybatis_Session {
			session = argInterface.(Session)
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
			if arg.Kind() != reflect.Ptr {
				paramMap[DefaultOneArg] = SqlArg{
					Value: argInterface,
					Type:  arg.Type(),
				}
			}
		}
	}
	result, err := sqlBuilder.BuildSql(paramMap, mapperXml, enableLog)
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

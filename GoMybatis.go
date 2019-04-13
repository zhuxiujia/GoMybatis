package GoMybatis

import (
	"github.com/zhuxiujia/GoMybatis/ast"
	"github.com/zhuxiujia/GoMybatis/lib/github.com/beevik/etree"
	"github.com/zhuxiujia/GoMybatis/utils"
	"bytes"
	"reflect"
	"strings"
)

const NewSessionFunc = "NewSession" //NewSession method,auto write implement body code

type Mapper struct {
	xml   *etree.Element
	nodes []ast.Node
}

//推荐默认使用单例传入
//根据sessionEngine写入到mapperPtr，value:指向mapper指针反射对象，xml：xml数据，sessionEngine：session引擎，enableLog:是否允许日志输出，log：日志实现
func WriteMapperByValue(value reflect.Value, xml []byte, sessionEngine SessionEngine) {
	if value.Kind() != reflect.Ptr {
		panic("UseMapper: UseMapper arg must be a pointer")
	}
	WriteMapper(value, xml, sessionEngine.SessionFactory(), sessionEngine.TempleteDecoder(), sessionEngine.SqlResultDecoder(), sessionEngine.SqlBuilder(), sessionEngine.LogEnable())
}

//推荐默认使用单例传入
//根据sessionEngine写入到mapperPtr，ptr:指向mapper指针，xml：xml数据，sessionEngine：session引擎，enableLog:是否允许日志输出，log：日志实现
func WriteMapperPtrByEngine(ptr interface{}, xml []byte, sessionEngine SessionEngine) {
	v := reflect.ValueOf(ptr)
	if v.Kind() != reflect.Ptr {
		panic("UseMapper: UseMapper arg must be a pointer")
	}
	WriteMapperByValue(v, xml, sessionEngine)
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
func WriteMapper(bean reflect.Value, xml []byte, sessionFactory *SessionFactory, templeteDecoder TempleteDecoder, decoder SqlResultDecoder, sqlBuilder SqlBuilder, enableLog bool) {
	beanCheck(bean)
	var mapperTree = LoadMapperXml(xml)
	templeteDecoder.DecodeTree(mapperTree, bean.Type())
	//构建期使用的map，无需考虑并发安全
	var methodXmlMap = makeMethodXmlMap(bean, mapperTree, sqlBuilder)
	var resultMaps = makeResultMaps(mapperTree)
	var returnTypeMap = makeReturnTypeMap(bean)
	var beanName = bean.Type().PkgPath() + bean.Type().String()

	UseMapperValue(bean, func(funcField reflect.StructField) func(args []reflect.Value, tagArgs []TagArg) []reflect.Value {
		//构建期
		var funcName = funcField.Name
		var returnType = returnTypeMap[funcName]
		if returnType == nil {
			panic("[GoMybatis] struct have no return values!")
		}
		//mapper
		var mapper = methodXmlMap[funcName]
		//resultMaps
		var resultMap map[string]*ResultProperty

		if funcName != NewSessionFunc {
			var resultMapId = mapper.xml.SelectAttrValue(Element_ResultMap, "")
			if resultMapId != "" {
				resultMap = resultMaps[resultMapId]
			}
		}

		//执行期
		if funcName == NewSessionFunc {
			var proxyFunc = func(args []reflect.Value, tagArgs []TagArg) []reflect.Value {
				var returnValue *reflect.Value = nil
				//build return Type
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
				var session Session
				var err error
				if session != nil {
					returnValue.Elem().Set(reflect.ValueOf(session).Elem().Addr().Convert(*returnType.ReturnOutType))
				} else {
					//err = utils.NewError("GoMybatis", "Create Session fail.arg session not exist!")
					session = sessionFactory.NewSession(beanName, SessionType_Default)
				}
				return buildReturnValues(returnType, returnValue, err)
			}
			return proxyFunc
		} else {
			var proxyFunc = func(args []reflect.Value, tagArgs []TagArg) []reflect.Value {
				var returnValue *reflect.Value = nil
				//build return Type
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
				//exe sql
				var e = exeMethodByXml(mapper.xml.Tag, beanName, sessionFactory, tagArgs, args, mapper.nodes, resultMap, returnValue, decoder, sqlBuilder, enableLog)
				return buildReturnValues(returnType, returnValue, e)
			}
			return proxyFunc
		}
	})
}

//check beans
func beanCheck(value reflect.Value) {
	var t = value.Type()
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	for i := 0; i < t.NumField(); i++ {
		var fieldItem = t.Field(i)
		if fieldItem.Type.Kind() != reflect.Func {
			continue
		}
		var argsLen = fieldItem.Type.NumIn() //参数长度，除session参数外。
		var customLen = 0
		for argIndex := 0; argIndex < fieldItem.Type.NumIn(); argIndex++ {
			var inType = fieldItem.Type.In(argIndex)
			if isCustomStruct(inType) {
				customLen++
			}
		}
		if argsLen > 1 && customLen > 1 {
			panic(`[GoMybats] ` + fieldItem.Name + ` must add tag "mapperParams:"*,*..."`)
		}
	}
}

func buildReturnValues(returnType *ReturnType, returnValue *reflect.Value, e error) []reflect.Value {
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

func makeReturnTypeMap(value reflect.Value) (returnMap map[string]*ReturnType) {
	returnMap = make(map[string]*ReturnType)
	var proxyType = value.Elem().Type()
	for i := 0; i < proxyType.NumField(); i++ {
		var funcType = proxyType.Field(i).Type
		var funcName = proxyType.Field(i).Name

		if funcType.Kind() != reflect.Func {
			continue
		}

		var numOut = funcType.NumOut()
		if numOut > 2 || numOut == 0 {
			panic("[GoMybatis] func '" + funcName + "()' return num out must = 1 or = 2!")
		}
		for f := 0; f < numOut; f++ {
			var outType = funcType.Out(f)
			if funcName != NewSessionFunc {
				//过滤NewSession方法
				if outType.Kind() == reflect.Ptr || (outType.Kind() == reflect.Interface && outType.String() != "error") {
					panic("[GoMybatis] func '" + funcName + "()' return '" + outType.String() + "' can not be a 'ptr' or 'interface'!")
				}
			}
			var returnType = returnMap[funcName]
			if returnType == nil {
				returnMap[funcName] = &ReturnType{
					ReturnIndex: -1,
					NumOut:      numOut,
				}
			}
			if outType.String() != "error" {
				returnMap[funcName].ReturnIndex = f
				returnMap[funcName].ReturnOutType = &outType
			} else {
				if returnMap[funcName].ErrorType != nil {
					panic("[GoMybatis] func '" + funcName + "()' must only return one 'error'!")
				}
				returnMap[funcName].ErrorType = &outType
			}
		}
	}
	return returnMap
}

//map[id]map[cloum]Property
func makeResultMaps(xmls map[string]etree.Token) map[string]map[string]*ResultProperty {
	var resultMaps = make(map[string]map[string]*ResultProperty)
	for _, item := range xmls {
		var typeString = reflect.TypeOf(item).String()
		if typeString == "*etree.Element" {
			var xmlItem = item.(*etree.Element)
			if xmlItem.Tag == Element_ResultMap {
				var resultPropertyMap = make(map[string]*ResultProperty)
				for _, elementItem := range xmlItem.ChildElements() {
					var property = ResultProperty{
						XMLName:  elementItem.Tag,
						Column:   elementItem.SelectAttrValue("column", ""),
						Property: elementItem.SelectAttrValue("property", ""),
						GoType:   elementItem.SelectAttrValue("goType", ""),
					}
					resultPropertyMap[property.Column] = &property
				}
				resultMaps[xmlItem.SelectAttrValue("id", "")] = resultPropertyMap
			}
		}
	}
	return resultMaps
}

//return a map map[`method`]*MapperXml
func makeMethodXmlMap(bean reflect.Value, mapperTree map[string]etree.Token, sqlBuilder SqlBuilder) map[string]*Mapper {
	var beanType = bean.Type()
	if beanType.Kind() == reflect.Ptr {
		beanType = beanType.Elem()
	}

	var methodXmlMap = make(map[string]*Mapper)
	var totalField = beanType.NumField()
	for i := 0; i < totalField; i++ {
		var fieldItem = beanType.Field(i)
		if fieldItem.Type.Kind() == reflect.Func {
			//field must be func
			methodFieldCheck(&beanType, &fieldItem)
			var mapperXml = findMapperXml(mapperTree, fieldItem.Name)
			if mapperXml != nil {
				methodXmlMap[fieldItem.Name] = &Mapper{
					xml:   mapperXml,
					nodes: sqlBuilder.NodeParser().ParserNodes(mapperXml.Child),
				}
			} else {
				if fieldItem.Name == NewSessionFunc {
					//过滤NewSession方法
					continue
				}
				panic("[GoMybatis] can not find method " + beanType.String() + "." + fieldItem.Name + "() in xml !")
			}
		}
	}
	return methodXmlMap
}

//方法基本规则检查
func methodFieldCheck(beanType *reflect.Type, methodType *reflect.StructField) {
	if methodType.Type.NumOut() < 1 {
		var buffer bytes.Buffer
		buffer.WriteString("[GoMybatis] bean ")
		buffer.WriteString((*beanType).Name())
		buffer.WriteString(".")
		buffer.WriteString(methodType.Name)
		buffer.WriteString("() must be return a 'error' type!")
		panic(buffer.String())
	}
	var errorTypeNum = 0
	for i := 0; i < methodType.Type.NumOut(); i++ {
		var outType = methodType.Type.Out(i)
		if outType.Kind() == reflect.Interface && outType.String() == "error" {
			errorTypeNum++
		}
	}
	if errorTypeNum != 1 {
		var buffer bytes.Buffer
		buffer.WriteString("[GoMybatis] bean ")
		buffer.WriteString((*beanType).Name())
		buffer.WriteString(".")
		buffer.WriteString(methodType.Name)
		buffer.WriteString("() must be return a 'error' type!")
		panic(buffer.String())
	}
}

func findMapperXml(mapperTree map[string]etree.Token, methodName string) *etree.Element {
	for _, mapperXml := range mapperTree {
		//exec sql,return data
		var typeString = reflect.TypeOf(mapperXml).String()
		if typeString == "*etree.Element" {
			var key = mapperXml.(*etree.Element).SelectAttrValue("id", "")
			if strings.EqualFold(key, methodName) {
				return mapperXml.(*etree.Element)
			}
		}
	}
	return nil
}

func exeMethodByXml(elementType ElementType, beanName string, sessionFactory *SessionFactory, tagParamMap []TagArg, args []reflect.Value, nodes []ast.Node, resultMap map[string]*ResultProperty, returnValue *reflect.Value, decoder SqlResultDecoder, sqlBuilder SqlBuilder, enableLog bool) error {
	//build sql string
	var session Session
	var sql string
	var err error
	session, sql, err = buildSql(tagParamMap, args, nodes, sqlBuilder, enableLog)
	if err != nil {
		return err
	}
	if sessionFactory == nil && session == nil {
		panic("[GoMybatis] exe sql need a SessionFactory or Session!")
	}
	//session
	if session == nil {
		session = sessionFactory.NewSession(beanName, SessionType_Default)
		//not arg session,just close!
		defer closeSession(sessionFactory, session)
	}
	var haveLastReturnValue = returnValue != nil && (*returnValue).IsNil() == false
	//do CRUD
	if elementType == Element_Select && haveLastReturnValue {
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

func buildSql(tagArgs []TagArg, args []reflect.Value, nodes []ast.Node, sqlBuilder SqlBuilder, enableLog bool) (Session, string, error) {
	var session Session
	var paramMap = make(map[string]interface{})
	var tagArgsLen = len(tagArgs)
	var argsLen = len(args) //参数长度，除session参数外。
	var customLen = 0
	var customIndex = -1
	for argIndex, arg := range args {
		var argInterface = arg.Interface()
		if arg.Kind() == reflect.Ptr && arg.IsNil() == false && argInterface != nil && arg.Type().String() == GoMybatis_Session_Ptr {
			session = *(argInterface.(*Session))
			continue
		} else if argInterface != nil && arg.Kind() == reflect.Interface && arg.Type().String() == GoMybatis_Session {
			session = argInterface.(Session)
			continue
		}
		if isCustomStruct(arg.Type()) {
			customLen++
			customIndex = argIndex
		}
		if arg.Type().String() == GoMybatis_Session_Ptr || arg.Type().String() == GoMybatis_Session {
			if argsLen > 0 {
				argsLen--
			}
			if tagArgsLen > 0 {
				tagArgsLen--
			}
		}
		if tagArgsLen > 0 && argIndex < tagArgsLen && tagArgs[argIndex].Name != "" {
			//插入2份参数，兼容大小写不敏感的参数
			var lowerKey = utils.LowerFieldFirstName(tagArgs[argIndex].Name)
			var upperKey = utils.UpperFieldFirstName(tagArgs[argIndex].Name)
			paramMap[lowerKey] = argInterface
			paramMap[upperKey] = argInterface
			//paramMap["type_"+lowerKey] = arg.Type()
			//paramMap["type_"+upperKey] = arg.Type()
		} else {
			paramMap[DefaultOneArg] = argInterface
			//paramMap["type_"+DefaultOneArg] = arg.Type()

		}
	}
	if customLen == 1 && customIndex != -1 {
		//只有一个结构体参数，需要展开它的成员变量 加入到map
		paramMap = scanStructArgFields(args[customIndex], nil)
	}

	result, err := sqlBuilder.BuildSql(paramMap, nodes)
	return session, result, err
}

//scan params
func scanStructArgFields(v reflect.Value, typeConvert func(arg interface{}) interface{}) map[string]interface{} {
	var t = v.Type()
	parameters := make(map[string]interface{})
	if v.Kind() == reflect.Ptr {
		if v.IsNil() == true {
			return parameters
		}
		//为指针，解引用
		v = v.Elem()
		t = t.Elem()
	}
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
			//parameters["type_"+jsonKey] = v.Field(i).Type()
		} else {
			parameters[typeValue.Name] = obj
			//parameters["type_"+typeValue.Name] = v.Field(i).Type()
		}
	}
	return parameters
}

func isCustomStruct(value reflect.Type) bool {
	if value.Kind() == reflect.Struct && value.String() != GoMybatis_Time && value.String() != GoMybatis_Time_Ptr {
		return true
	} else {
		return false
	}
}

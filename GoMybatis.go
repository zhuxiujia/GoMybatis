package GoMybatis

import (
	"bytes"
	"github.com/zhuxiujia/GoMybatis/ast"
	"github.com/zhuxiujia/GoMybatis/lib/github.com/beevik/etree"
	"github.com/zhuxiujia/GoMybatis/utils"
	"log"
	"reflect"
	"strconv"
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
		panic("AopProxy: AopProxy arg must be a pointer")
	}
	WriteMapper(value, xml, sessionEngine)
	sessionEngine.RegisterObj(value.Interface(), value.Type().Elem().Name())
}

//推荐默认使用单例传入
//根据sessionEngine写入到mapperPtr，ptr:指向mapper指针，xml：xml数据，sessionEngine：session引擎，enableLog:是否允许日志输出，log：日志实现
func WriteMapperPtrByEngine(ptr interface{}, xml []byte, sessionEngine SessionEngine) {
	v := reflect.ValueOf(ptr)
	if v.Kind() != reflect.Ptr {
		panic("AopProxy: AopProxy arg must be a pointer")
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
func WriteMapper(bean reflect.Value, xml []byte, sessionEngine SessionEngine) {
	beanCheck(bean)
	var mapperTree = LoadMapperXml(xml)
	sessionEngine.TempleteDecoder().DecodeTree(mapperTree, bean.Type())
	//构建期使用的map，无需考虑并发安全
	var methodXmlMap = makeMethodXmlMap(bean, mapperTree, sessionEngine.SqlBuilder())
	var resultMaps = makeResultMaps(mapperTree)
	var returnTypeMap = makeReturnTypeMap(bean.Elem().Type())
	var beanName = bean.Type().PkgPath() + bean.Type().String()

	ProxyValue(bean, func(funcField reflect.StructField, field reflect.Value) func(arg ProxyArg) []reflect.Value {
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
			var proxyFunc = func(arg ProxyArg) []reflect.Value {
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
				var session = sessionEngine.SessionFactory().NewSession(beanName, SessionType_Default)
				var err error
				returnValue.Elem().Set(reflect.ValueOf(session).Elem().Addr().Convert(*returnType.ReturnOutType))
				return buildReturnValues(returnType, returnValue, err)
			}
			return proxyFunc
		} else {
			var proxyFunc = func(arg ProxyArg) []reflect.Value {
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
				var e = exeMethodByXml(mapper.xml.Tag, beanName, sessionEngine, arg, mapper.nodes, resultMap, returnValue)
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

func makeReturnTypeMap(value reflect.Type) (returnMap map[string]*ReturnType) {
	returnMap = make(map[string]*ReturnType)
	var proxyType = value
	for i := 0; i < proxyType.NumField(); i++ {
		var funcType = proxyType.Field(i).Type
		var funcName = proxyType.Field(i).Name

		if funcType.Kind() != reflect.Func {
			if funcType.Kind() == reflect.Struct {
				var childMap = makeReturnTypeMap(funcType)
				for k, v := range childMap {
					returnMap[k] = v
				}
			}
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
				//error
				returnMap[funcName].ErrorType = &outType
			}
		}
		if returnMap[funcName].ErrorType == nil {
			panic("[GoMybatis] func '" + funcName + "()' must return an 'error'!")
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
						LangType: elementItem.SelectAttrValue("langType", ""),
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
					nodes: sqlBuilder.NodeParser().Parser(mapperXml.Child),
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

	var mapperParams = methodType.Tag.Get("mapperParams")
	if methodType.Type.NumOut() > 1 && mapperParams == "" && !(methodType.Name == "NewSession") {
		log.Println("[GoMybatis] warning ======================== " + (*beanType).Name() + "." + methodType.Name + "() have not define tag mapperParams:\"\",maybe can not get param value!")
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

func exeMethodByXml(elementType ElementType, beanName string, sessionEngine SessionEngine, proxyArg ProxyArg, nodes []ast.Node, resultMap map[string]*ResultProperty, returnValue *reflect.Value) error {
	//TODO　CallBack and Session must Location in build step!
	var session Session
	var sql string
	var err error
	var array_arg = []interface{}{}
	session, sql, err = buildSql(proxyArg, nodes, sessionEngine.SqlBuilder(), &array_arg)
	if err != nil {
		return err
	}
	if sessionEngine.SessionFactory() == nil && session == nil {
		panic("[GoMybatis] exe sql need a SessionFactory or Session!")
	}
	//session
	if session == nil {
		var goroutineID int64 //协程id
		if sessionEngine.GoroutineIDEnable() {
			goroutineID = utils.GoroutineID()
		} else {
			goroutineID = 0
		}
		session = sessionEngine.GoroutineSessionMap().Get(goroutineID)
	}
	if session == nil {
		var s, err = sessionEngine.NewSession(beanName)
		if err != nil {
			return err
		}
		session = s
		defer session.Close()
	}
	var haveLastReturnValue = returnValue != nil && (*returnValue).IsNil() == false
	//do CRUD
	if elementType == Element_Select && haveLastReturnValue {
		//is select and have return value
		if sessionEngine.LogEnable() {
			sessionEngine.LogSystem().SendLog("[GoMybatis] [", session.Id(), "] Query ==> "+sql)
			sessionEngine.LogSystem().SendLog("[GoMybatis] [", session.Id(), "] Args  ==> "+utils.SprintArray(array_arg))
		}
		res, err := session.QueryPrepare(sql, array_arg...)
		defer func() {
			if sessionEngine.LogEnable() {
				var RowsAffected = "0"
				if err == nil && res != nil {
					RowsAffected = strconv.Itoa(len(res))
				}
				sessionEngine.LogSystem().SendLog("[GoMybatis] [", session.Id(), "] ReturnRows <== "+RowsAffected)
				if err != nil {
					sessionEngine.LogSystem().SendLog("[GoMybatis] [", session.Id(), "] error == "+err.Error())
				}
			}
		}()
		if err != nil {
			return err
		}
		err = sessionEngine.SqlResultDecoder().Decode(resultMap, res, returnValue.Interface())
		if err != nil {
			return err
		}
	} else {
		if sessionEngine.LogEnable() {
			sessionEngine.LogSystem().SendLog("[GoMybatis] [", session.Id(), "] Exec ==> "+sql)
			sessionEngine.LogSystem().SendLog("[GoMybatis] [", session.Id(), "] Args ==> "+utils.SprintArray(array_arg))
		}
		var res, err = session.ExecPrepare(sql, array_arg...)
		defer func() {
			if sessionEngine.LogEnable() {
				var RowsAffected = "0"
				if err == nil && res != nil {
					RowsAffected = strconv.FormatInt(res.RowsAffected, 10)
				}
				sessionEngine.LogSystem().SendLog("[GoMybatis] [", session.Id(), "] RowsAffected <== "+RowsAffected)
				if err != nil {
					sessionEngine.LogSystem().SendLog("[GoMybatis] [", session.Id(), "] error == "+err.Error())
				}
			}
		}()

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

func buildSql(proxyArg ProxyArg, nodes []ast.Node, sqlBuilder SqlBuilder, array_arg *[]interface{}) (Session, string, error) {
	var session Session
	var paramMap = make(map[string]interface{})
	var tagArgsLen = proxyArg.TagArgsLen
	var argsLen = proxyArg.ArgsLen //参数长度，除session参数外。
	var customLen = 0
	var customIndex = -1
	for argIndex, arg := range proxyArg.Args {
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
		if tagArgsLen > 0 && argIndex < tagArgsLen && proxyArg.TagArgs[argIndex].Name != "" {
			//插入2份参数，兼容大小写不敏感的参数
			var lowerKey = utils.LowerFieldFirstName(proxyArg.TagArgs[argIndex].Name)
			var upperKey = utils.UpperFieldFirstName(proxyArg.TagArgs[argIndex].Name)
			paramMap[lowerKey] = argInterface
			paramMap[upperKey] = argInterface
		} else {
			//未命名参数，为arg加参数位置，例如 arg0,arg1,arg2....
			paramMap[DefaultOneArg+strconv.Itoa(argIndex)] = argInterface
		}
	}
	if customLen == 1 && customIndex != -1 {
		//只有一个结构体参数，需要展开它的成员变量 加入到map
		var tag *TagArg
		if proxyArg.TagArgsLen == 1 {
			tag = &proxyArg.TagArgs[0]
		}
		paramMap = scanStructArgFields(proxyArg.Args[customIndex], tag)
	}

	result, err := sqlBuilder.BuildSql(paramMap, nodes, array_arg)
	return session, result, err
}

//scan params
func scanStructArgFields(v reflect.Value, tag *TagArg) map[string]interface{} {
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

	var structArg = make(map[string]interface{})

	//json arg,性能较差
	//var vptr=v.Interface()
	//var js,_=json.Marshal(vptr)
	//json.Unmarshal(js,&structArg)
	//
	//for key,value:=range structArg {
	//	parameters[key]=value
	//}

	//reflect arg,性能较快
	for i := 0; i < t.NumField(); i++ {
		var typeValue = t.Field(i)
		var field = v.Field(i)

		var obj interface{}
		if field.CanInterface() {
			obj = field.Interface()
		}
		var jsonKey = typeValue.Tag.Get(`json`)
		if strings.Index(jsonKey, ",") != -1 {
			jsonKey = strings.Split(jsonKey, ",")[0]
		}
		if jsonKey != "" {
			parameters[jsonKey] = obj
			structArg[jsonKey] = obj
			parameters[typeValue.Name] = obj
			structArg[typeValue.Name] = obj
		} else {
			parameters[typeValue.Name] = obj
			structArg[typeValue.Name] = obj
		}
	}
	if tag != nil && parameters[tag.Name] == nil {
		parameters[tag.Name] = structArg
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

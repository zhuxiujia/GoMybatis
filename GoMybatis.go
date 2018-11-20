package GoMybatis

import (
	"strings"
	"reflect"
	"errors"
	"fmt"
)

const SessionIdKey = "sessionId"

//如果使用UseProxyMapperByEngine，则内建默认的SessionFactory
var DefaultSessionFactory SessionFactory

//bean 工厂，根据xml配置创建函数,并且动态代理到你定义的struct func里
//bean 参数必须为指针类型,指向你定义的struct
//你定义的struct必须有可导出的func属性,例如：
//type MyUserMapperImpl struct {
//	UserMapper                                                 `mapperPath:"/mapper/user/UserMapper.xml"`
//	SelectById    func(id string, result *model.User) error    `mapperParams:"id"`
//	SelectByPhone func(id string, phone string, result *model.User) error `mapperParams:"id,phone"`
//	DeleteById    func(id string, result *int64) error         `mapperParams:"id"`
//	Insert        func(arg model.User, result *int64) error
//}
//func的参数支持2种函数，第一种函数 基本参数个数无限制(并且需要用Tag指定参数名逗号隔开,例如`mapperParams:"id,phone"`)，最后一个参数必须为返回数据类型的指针(例如result *model.User)，返回值为error
//func的参数支持2种函数，第二种函数第一个参数必须为结构体(例如 arg model.User,该结构体的属性可以指定tag `json:"xxx"`为参数名称),最后一个参数必须为返回数据类型的指针(例如result *model.User)，返回值为error
//使用UseProxyMapper函数设置代理后即可正常使用。
func UseProxyMapperByEngine(bean interface{}, xml []byte, sqlEngine *SessionEngine) {
	v := reflect.ValueOf(bean)
	if v.Kind() != reflect.Ptr {
		panic("UseMapper: UseMapper arg must be a pointer")
	}
	DefaultSessionFactory = SessionFactory{}.New(sqlEngine)
	UseProxyMapperFromValue(v, xml, &DefaultSessionFactory)
}

//bean 工厂，根据xml配置创建函数,并且动态代理到你定义的struct func里
//bean 参数必须为指针类型,指向你定义的struct
//你定义的struct必须有可导出的func属性,例如：
//type MyUserMapperImpl struct {
//	UserMapper                                                 `mapperPath:"/mapper/user/UserMapper.xml"`
//	SelectById    func(id string, result *model.User) error    `mapperParams:"id"`
//	SelectByPhone func(id string, phone string, result *model.User) error `mapperParams:"id,phone"`
//	DeleteById    func(id string, result *int64) error         `mapperParams:"id"`
//	Insert        func(arg model.User, result *int64) error
//}
//func的参数支持2种函数，第一种函数 基本参数个数无限制(并且需要用Tag指定参数名逗号隔开,例如`mapperParams:"id,phone"`)，最后一个参数必须为返回数据类型的指针(例如result *model.User)，返回值为error
//func的参数支持2种函数，第二种函数第一个参数必须为结构体(例如 arg model.User,该结构体的属性可以指定tag `json:"xxx"`为参数名称),最后一个参数必须为返回数据类型的指针(例如result *model.User)，返回值为error
//使用UseProxyMapper函数设置代理后即可正常使用。
func UseProxyMapperByFactory(bean interface{}, xml []byte, sqlEngine *SessionFactory) {
	v := reflect.ValueOf(bean)
	if v.Kind() != reflect.Ptr {
		panic("UseMapper: UseMapper arg must be a pointer")
	}
	UseProxyMapperFromValue(v, xml, sqlEngine)
}

//bean 工厂，根据xml配置创建函数,并且动态代理到你定义的struct func里
//bean 参数必须为reflect.Value
//你定义的struct必须有可导出的func属性,例如：
//type MyUserMapperImpl struct {
//	UserMapper                                                 `mapperPath:"/mapper/user/UserMapper.xml"`
//	SelectById    func(id string, result *model.User) error    `mapperParams:"id"`
//	SelectByPhone func(id string, phone string, result *model.User) error `mapperParams:"id,phone"`
//	DeleteById    func(id string, result *int64) error         `mapperParams:"id"`
//	Insert        func(arg model.User, result *int64) error
//}
//func的参数支持2种函数，第一种函数 基本参数个数无限制(并且需要用Tag指定参数名逗号隔开,例如`mapperParams:"id,phone"`)，最后一个参数必须为返回数据类型的指针(例如result *model.User)，返回值为error
//func的参数支持2种函数，第二种函数第一个参数必须为结构体(例如 arg model.User,该结构体的属性可以指定tag `json:"xxx"`为参数名称),最后一个参数必须为返回数据类型的指针(例如result *model.User)，返回值为error
//使用UseProxyMapper函数设置代理后即可正常使用。
func UseProxyMapperFromValue(bean reflect.Value, xml []byte, sessionFactory *SessionFactory) {
	var mapperTree = LoadMapperXml(xml)
	var proxyFunc = func(method string, args []reflect.Value, tagParams []string) error {
		var lastArgsIndex = len(args) - 1
		var tagParamsLen = len(tagParams)
		var argsLen = len(args)
		var lastArgValue *reflect.Value = nil
		if argsLen != 0 && args[lastArgsIndex].Kind() == reflect.Ptr {
			lastArgValue = &args[lastArgsIndex]
			if lastArgValue.Kind() != reflect.Ptr {
				//最后一个参数必须为指针，或者不传任何参数
				return errors.New(`[method params last param must be pointer!],method =` + method)
			}
		}
		//build params
		var paramMap = make(map[string]interface{})
		if tagParamsLen != 0 {
			for index, v := range tagParams {
				paramMap[v] = args[index].Interface()
			}
		}
		var findMethod = false
		for _, mapperXml := range mapperTree {
			//exec sql,return data
			if strings.EqualFold(mapperXml.Id, method) {
				findMethod = true
				//build sql string
				var sql string
				var err error
				sql, err = buildSql(paramMap, args, mapperXml)
				if err != nil {
					return err
				}
				//session
				var session *Session
				var sessionId = paramMap[SessionIdKey]
				if sessionId != nil {
					session = sessionFactory.GetSession(sessionId.(string))
					if session == nil {
						return errors.New("[GoMybatis] session is nil! sessionId=" + fmt.Sprint(sessionId))
					}
				} else {
					session = sessionFactory.NewSession()
				}
				var haveLastReturnValue = lastArgValue != nil && (*lastArgValue).IsNil() == false
				//do CRUD
				if mapperXml.Tag == Select && haveLastReturnValue {
					//is select and have return value
					results, err := (*session).Query(sql)
					if err != nil {
						return err
					}
					err = Unmarshal(results, lastArgValue.Interface())
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
				if sessionId == nil {
					sessionFactory.CloseSession((*session).Id())
				}
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

func buildSql(paramMap map[string]interface{}, args []reflect.Value, mapperXml MapperXml) (string, error) {
	for _, arg := range args {
		if arg.Kind() == reflect.Ptr {
			//指针，则退出
			continue
		}
		if arg.Kind() == reflect.Struct && arg.Type().String() != `time.Time` {
			paramMap = scanParamterBean(arg.Interface(), nil)
		} else {
			paramMap[DefaultOneArg] = arg.Interface()
		}
	}
	return BuildSqlFromMap(paramMap, mapperXml)
}

//scan params
func scanParamterBean(arg interface{}, typeConvert func(arg interface{}) interface{}) map[string]interface{} {
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

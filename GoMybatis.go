package GoMybatis

import (
	"encoding/xml"
	"github.com/beevik/etree"
	"strings"
	"reflect"
	"fmt"
	"errors"
	"regexp"
	"github.com/go-xorm/xorm"
	"github.com/zhuxiujia/GoMybatis/lib/govaluate"
)

type Mapper struct {
	XMLName    xml.Name    `xml:"mapper"`
	ResultMaps []ResultMap `xml:"resultMap"`
}

type MapperXml struct {
	Tag          string
	Id           string
	ElementItems []ElementItem
}

type ElementItem struct {
	ElementType  string
	Propertys    map[string]string
	DataString   string
	ElementItems []ElementItem
}

const DefaultOneArg = `[default]`

//读取xml
func LoadMapperXml(bytes []byte) (items []MapperXml) {
	var m Mapper
	xml.Unmarshal(bytes, &m)

	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(bytes); err != nil {
		panic(err)
	}
	root := doc.SelectElement("mapper")
	for _, s := range root.ChildElements() {
		var attrMap = attrToProperty(s.Attr)
		var elItems = loop(s)
		if s.Tag == Insert || s.Tag == Delete || s.Tag == Update || s.Tag == Select {
			items = append(items, MapperXml{
				Tag:          s.Tag,
				Id:           attrMap[`id`],
				ElementItems: elItems,
			})
		}
	}
	return items
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
func UseProxyMapper(bean interface{}, xml []byte, xormEngine *xorm.Engine) {
	v := reflect.ValueOf(bean)
	if v.Kind() != reflect.Ptr {
		panic("UseMapper: UseMapper arg must be a pointer")
	}
	UseProxyMapperFromValue(v, xml, xormEngine)
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
func UseProxyMapperFromValue(bean reflect.Value, xml []byte, engine *xorm.Engine) {
	var mapperTree = LoadMapperXml(xml)
	var proxyFunc = func(method string, args []reflect.Value, params []string) error {
		var paramsLen = len(params)
		var argsLen = len(args)
		var lastArgValue *reflect.Value = nil
		if argsLen != 0 && args[len(args)-1].Kind() == reflect.Ptr {
			lastArgValue = &args[len(args)-1]
			if lastArgValue.Kind() != reflect.Ptr {
				//最后一个参数必须为指针，或者不传任何参数
				return errors.New(`[method params last param must be pointer!],method =` + method)
			}
		}

		//build params
		var paramMap = make(map[string]interface{})
		if paramsLen != 0 {
			for index, v := range params {
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
				if paramsLen != 0 {
					sql, err = buildSqlFromMap(paramMap, mapperXml)
				} else if paramsLen == 0 && argsLen == 0 {
					sql, err = buildSqlFromMap(paramMap, mapperXml)
				} else {
					sql, err = buildSql(args[0], mapperXml)
				}
				if err != nil {
					return err
				}
				//TODO do CRUD
				if mapperXml.Tag == Select {
					if lastArgValue != nil && (*lastArgValue).IsNil() == false {
						results, err := engine.Query(sql)
						if err != nil {
							return err
						}
						err = Unmarshal(results, lastArgValue.Interface())
						if err != nil {
							return err
						}
					} else {
						var _, err = engine.Exec(sql)
						if err != nil {
							return err
						}
					}
				} else if mapperXml.Tag == Update || mapperXml.Tag == Delete || mapperXml.Tag == Insert {
					var res, err = engine.Exec(sql)
					if lastArgValue != nil {
						if lastArgValue.IsNil() == false {
							var rows, err = res.RowsAffected()
							if err != nil {
								return err
							}
							lastArgValue.Elem().Set(reflect.ValueOf(rows))
						}
					}
					if err != nil {
						return err
					}
				}
				//匹配完成退出
				break
			}
		}
		if findMethod == false {
			return errors.New(`[not method find at xml file],method =` + method)
		}
		return nil
	}
	UseMapperValue(bean, proxyFunc)
}

func buildSqlFromMap(paramMap map[string]interface{}, mapperXml MapperXml) (string, error) {
	var sql = ``
	sql, err := createFromElement(mapperXml.ElementItems, sql, paramMap)
	if err != nil {
		return sql, err
	}
	fmt.Println("[create sql]", sql)
	return sql, nil
}

func buildSql(arg0 reflect.Value, mapperXml MapperXml) (string, error) {
	var params = make(map[string]interface{})
	if arg0.Kind() == reflect.Struct && arg0.Type().String() != `time.Time` {
		params = scanParamterBean(arg0.Interface(), nil)
	} else {
		params[DefaultOneArg] = arg0.Interface()
	}
	return buildSqlFromMap(params, mapperXml)
}

func attrToProperty(attrs []etree.Attr) map[string]string {
	var m = make(map[string]string)
	for _, v := range attrs {
		m[v.Key] = v.Value
	}
	return m
}

func loop(element *etree.Element) []ElementItem {
	var els = make([]ElementItem, 0)
	for _, el := range element.Child {
		var typeString = reflect.ValueOf(el).Type().String()
		if typeString == `*etree.CharData` {
			var d = el.(*etree.CharData)
			var str = d.Data
			str = ReplaceAllBlankSpace(str)
			if str != "" {
				var elementItem = ElementItem{
					ElementType: String,
					DataString:  str,
				}
				els = append(els, elementItem)
			}
		} else if typeString == `*etree.Element` {
			var e = el.(*etree.Element)
			var element = ElementItem{
				ElementType:  e.Tag,
				ElementItems: make([]ElementItem, 0),
				Propertys:    attrToProperty(e.Attr),
			}
			if len(e.Child) > 0 {
				var loopEls = loop(e)
				for _, item := range loopEls {
					element.ElementItems = append(element.ElementItems, item)
				}
			}
			els = append(els, element)
		}
	}
	return els
}

func createFromElement(itemTree []ElementItem, sql string, param map[string]interface{}) (result string, err error) {
	for _, v := range itemTree {
		var loopChildItem = true
		if v.ElementType == String {
			//string element
			sql += repleaceArg(v.DataString, param, DefaultSqlTypeConvertFunc)
		} else if v.ElementType == If {
			//if element
			var test = v.Propertys[`test`]
			var andStrings = strings.Split(test, ` and `)
			for index, expression := range andStrings {
				//test表达式解析
				var evaluateParameters = scanParamterMap(param, DefaultExpressionTypeConvertFunc)
				expression = expressionToIfZeroExpression(evaluateParameters, expression)
				evalExpression, err := govaluate.NewEvaluableExpression(expression)
				if err != nil {
					fmt.Println(err)
				}
				result, err := evalExpression.Evaluate(evaluateParameters)
				if err != nil {
					err = errors.New("test() -> `" + expression + "` " + err.Error())
					return sql, err
				}
				if result.(bool) {
					//test表达式成立
					if index == (len(andStrings) - 1) {
						sql += repleaceArg(v.DataString, param, DefaultSqlTypeConvertFunc)
					}
				} else {
					loopChildItem = false
					break
				}
			}
		} else if v.ElementType == Trim {
			var prefix = v.Propertys[`prefix`]
			var suffix = v.Propertys[`suffix`]
			var suffixOverrides = v.Propertys[`suffixOverrides`]
			if v.ElementItems != nil && len(v.ElementItems) > 0 && loopChildItem {
				var tempTrimSql = ``
				tempTrimSql, err = createFromElement(v.ElementItems, tempTrimSql, param)
				if err != nil {
					return tempTrimSql, err
				}
				tempTrimSql = strings.Trim(tempTrimSql, ` `)
				tempTrimSql = strings.Trim(tempTrimSql, suffixOverrides)
				tempTrimSql = prefix + ` ` + tempTrimSql + ` ` + suffix
				sql = sql + tempTrimSql
				loopChildItem = false
			}
		} else if v.ElementType == Set {
			var suffixOverrides = ","
			if v.ElementItems != nil && len(v.ElementItems) > 0 && loopChildItem {
				var trim = ``
				trim, err = createFromElement(v.ElementItems, trim, param)
				if err != nil {
					return trim, err
				}
				trim = strings.Trim(trim, ` `)
				trim = strings.Trim(trim, suffixOverrides)
				trim = ` set ` + trim + ` `
				sql = sql + trim
				loopChildItem = false
			}
		} else if v.ElementType == Foreach {
			var collection = v.Propertys[`collection`]
			var index = v.Propertys[`index`]
			var item = v.Propertys[`item`]
			var open = v.Propertys[`open`]
			var close = v.Propertys[`close`]
			var separator = v.Propertys[`separator`]
			var tempSql = ``
			var datas = param[collection]
			var collectionValue = reflect.ValueOf(datas)
			if collectionValue.Len() > 0 {
				for i := 0; i < collectionValue.Len(); i++ {
					var dataItem = collectionValue.Index(i).Interface()
					var tempParam = make(map[string]interface{})
					tempParam[item] = dataItem
					tempParam[index] = index
					for k, v := range param {
						tempParam[k] = v
					}
					if v.ElementItems != nil && len(v.ElementItems) > 0 && loopChildItem {
						tempSql, err = createFromElement(v.ElementItems, tempSql, tempParam)
						if err != nil {
							return tempSql, err
						}
					}
				}
			}
			tempSql = open + ` ` + tempSql + ` ` + close
			tempSql = strings.Trim(tempSql, separator)
			sql += tempSql
			loopChildItem = false
		}
		if v.ElementItems != nil && len(v.ElementItems) > 0 && loopChildItem {
			sql, err = createFromElement(v.ElementItems, sql, param)
			if err != nil {
				return sql, err
			}
		}
	}
	return sql, nil
}

//表达式 ''转换为 0
func expressionToIfZeroExpression(evaluateParameters map[string]interface{}, expression string) string {
	for k, v := range evaluateParameters {
		if strings.Index(expression, k) != -1 {
			var t = reflect.TypeOf(v)
			if t.String() != `string` {
				expression = strings.Replace(expression, `''`, `0`, -1)
			}
			return expression
		}
	}
	return expression
}

//替换参数
func repleaceArg(data string, parameters map[string]interface{}, typeConvertFunc func(arg interface{}) string) string {
	if data == "" {
		return data
	}
	for k, v := range parameters {
		if k == DefaultOneArg {
			var str = typeConvertFunc(v)
			re, _ := regexp.Compile("\\#\\{[^}]*\\}")
			data = re.ReplaceAllString(data, str)
		} else {
			var str = typeConvertFunc(v)
			re, _ := regexp.Compile("\\#\\{" + k + "[^}]*\\}")
			data = re.ReplaceAllString(data, str)
		}
	}
	return data
}

//scan params
func scanParamterMap(parameters map[string]interface{}, typeConvert func(arg interface{}) interface{}) map[string]interface{} {
	var newMap = make(map[string]interface{})
	for k, obj := range parameters {
		if typeConvert != nil {
			obj = typeConvert(obj)
		}
		newMap[k] = obj
	}
	return newMap
}

//scan params
func scanParamterBean(arg interface{}, typeConvert func(arg interface{}) interface{}) map[string]interface{} {
	parameters := make(map[string]interface{})
	v := reflect.ValueOf(arg)
	t := reflect.TypeOf(arg)
	if t.Kind() != reflect.Struct {
		panic(`the scanParamterBean() arg is not a struct type!,type =` + t.String())
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

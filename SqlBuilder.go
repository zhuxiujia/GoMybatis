package GoMybatis

import (
	"bytes"
	"fmt"
	"github.com/zhuxiujia/GoMybatis/utils"
	"log"
	"reflect"
	"strings"
)

type SqlBuilder interface {
	BuildSql(paramMap map[string]SqlArg, mapperXml *MapperXml, enableLog bool) (string, error)
	ExpressionEngineProxy() ExpressionEngineProxy
	SqlArgTypeConvert() SqlArgTypeConvert
	ExpressionTypeConvert() ExpressionTypeConvert
}

type GoMybatisSqlBuilder struct {
	expressionTypeConvert ExpressionTypeConvert
	sqlArgTypeConvert     SqlArgTypeConvert
	expressionEngineProxy      ExpressionEngineProxy
}

func (this GoMybatisSqlBuilder) ExpressionEngineProxy() ExpressionEngineProxy {
	return this.expressionEngineProxy
}
func (this GoMybatisSqlBuilder) SqlArgTypeConvert() SqlArgTypeConvert {
	return this.sqlArgTypeConvert
}
func (this GoMybatisSqlBuilder) ExpressionTypeConvert() ExpressionTypeConvert {
	return this.expressionTypeConvert
}

func (this GoMybatisSqlBuilder) New(ExpressionTypeConvert ExpressionTypeConvert, SqlArgTypeConvert SqlArgTypeConvert, expressionEngine ExpressionEngineProxy) GoMybatisSqlBuilder {
	this.expressionTypeConvert = ExpressionTypeConvert
	this.sqlArgTypeConvert = SqlArgTypeConvert
	this.expressionEngineProxy = expressionEngine
	return this
}

func (this GoMybatisSqlBuilder) BuildSql(paramMap map[string]SqlArg, mapperXml *MapperXml, enableLog bool) (string, error) {
	var sql bytes.Buffer
	err := this.createFromElement(mapperXml.ElementItems, &sql, paramMap)
	if err != nil {
		return "", err
	}
	var sqlStr = sql.String()
	sql.Reset()
	if enableLog {
		log.Println("[GoMybatis] Preparing sql ==> ", sqlStr)
	}
	return sqlStr, nil
}

func (this *GoMybatisSqlBuilder) createFromElement(itemTree []ElementItem, sql *bytes.Buffer, sqlArgMap map[string]SqlArg) error {
	if this.sqlArgTypeConvert == nil || this.expressionTypeConvert == nil {
		panic("[GoMybatis] GoMybatisSqlBuilder.sqlArgTypeConvert and GoMybatisSqlBuilder.expressionTypeConvert can not be nil!")
	}
	//默认的map[string]interface{}
	var defaultArgMap = this.makeArgInterfaceMap(sqlArgMap)
	//test表达式参数map
	var evaluateParameters map[string]interface{}
	for _, v := range itemTree {
		var loopChildItem = true
		var breakChildItem = false
		switch v.ElementType {
		case Element_bind:
			//bind,param args change!need update
			sqlArgMap = this.bindBindElementArg(sqlArgMap, v, this.sqlArgTypeConvert)
			defaultArgMap = this.makeArgInterfaceMap(sqlArgMap)
			if evaluateParameters != nil {
				evaluateParameters = this.expressParamterMap(sqlArgMap, this.expressionTypeConvert)
			}
			break
		case Element_String:
			//string element
			var replaceSql, err = replaceArg(v.DataString, defaultArgMap, this.sqlArgTypeConvert, &this.expressionEngineProxy)
			if err != nil {
				return err
			}
			sql.WriteString(replaceSql)
			break
		case Element_If:
			//if element
			var expression = v.Propertys[`test`]
			var result, err = this.doIfElement(expression, sqlArgMap, evaluateParameters)
			if err != nil {
				return err
			}
			if result {
				//test > true,write sql string
				var replaceSql, err = replaceArg(v.DataString, defaultArgMap, this.sqlArgTypeConvert, &this.expressionEngineProxy)
				if err != nil {
					return err
				}
				sql.WriteString(replaceSql)
			} else {
				// test > fail ,end loop
				loopChildItem = false
				break
			}
			break
		case Element_Trim:
			var prefix = v.Propertys[`prefix`]
			var suffix = v.Propertys[`suffix`]
			var suffixOverrides = v.Propertys[`suffixOverrides`]
			var prefixOverrides = v.Propertys[`prefixOverrides`]
			var err = this.elementTrim(&loopChildItem, v.ElementItems, sqlArgMap, prefix, suffix, prefixOverrides, suffixOverrides, sql)
			if err != nil {
				return err
			}
			break
		case Element_Set:
			if loopChildItem && v.ElementItems != nil && len(v.ElementItems) > 0 {
				var trim bytes.Buffer
				var err = this.createFromElement(v.ElementItems, &trim, sqlArgMap)
				if err != nil {
					return err
				}
				var trimString = strings.Trim(strings.Trim(trim.String(), " "), DefaultOverrides)
				trim.Reset()
				trim.WriteString(` `)
				trim.WriteString(` set `)
				trim.WriteString(trimString)
				trim.WriteString(` `)
				sql.Write(trim.Bytes())
				loopChildItem = false
			}
			break
		case Element_Foreach:
			var collection = v.Propertys[`collection`]
			var index = v.Propertys[`index`]
			var item = v.Propertys[`item`]
			var open = v.Propertys[`open`]
			var close = v.Propertys[`close`]
			var separator = v.Propertys[`separator`]
			if collection == "" {
				panic(`[GoMybatis] collection value can not be "" in <foreach collection=""> !`)
			}
			var tempSql bytes.Buffer
			var datas = sqlArgMap[collection].Value
			var collectionValue = reflect.ValueOf(datas)
			if collectionValue.Kind() != reflect.Slice && collectionValue.Kind() != reflect.Map {
				panic(`[GoMybatis] collection value must be a slice or map !`)
			}
			var collectionValueLen = collectionValue.Len()
			if collectionValueLen == 0 {
				continue
			}
			switch collectionValue.Kind() {
			case reflect.Map:
				var mapKeys = collectionValue.MapKeys()
				var collectionKeyLen = len(mapKeys)
				if collectionKeyLen == 0 {
					continue
				}
				for _, keyValue := range mapKeys {
					var key = keyValue.Interface()
					var collectionItem = collectionValue.MapIndex(keyValue)
					var tempArgMap = make(map[string]SqlArg) //temp parameter Map
					for k, v := range sqlArgMap {
						tempArgMap[k] = v
					}
					if item != "" {
						tempArgMap[item] = SqlArg{
							Value: collectionItem.Interface(),
							Type:  collectionItem.Type(),
						}
					}
					tempArgMap[index] = SqlArg{
						Value: key,
						Type:  keyValue.Type(),
					}
					if loopChildItem && v.ElementItems != nil && len(v.ElementItems) > 0 {
						var err = this.createFromElement(v.ElementItems, &tempSql, tempArgMap)
						if err != nil {
							return err
						}
						tempSql.WriteString(separator)
					}
				}
				break
			case reflect.Slice:
				for i := 0; i < collectionValueLen; i++ {
					var collectionItem = collectionValue.Index(i)
					var tempArgMap = make(map[string]SqlArg) //temp parameter Map
					for k, v := range sqlArgMap {
						tempArgMap[k] = v
					}
					if item != "" {
						tempArgMap[item] = SqlArg{
							Value: collectionItem.Interface(),
							Type:  collectionItem.Type(),
						}
					}
					if index != "" {
						tempArgMap[index] = SqlArg{
							Value: index,
							Type:  IntType,
						}
					}
					if loopChildItem && v.ElementItems != nil && len(v.ElementItems) > 0 {
						var err = this.createFromElement(v.ElementItems, &tempSql, tempArgMap)
						if err != nil {
							return err
						}
						tempSql.WriteString(separator)
					}
				}
				break
			}
			var newTempSql bytes.Buffer
			var tempSqlString = strings.Trim(strings.Trim(tempSql.String(), " "), separator)
			newTempSql.WriteString(open)
			newTempSql.WriteString(tempSqlString)
			newTempSql.WriteString(close)

			tempSql.Reset()
			sql.Write(newTempSql.Bytes())
			loopChildItem = false
			break
		case Element_choose:
			//read when and otherwise
			var temp bytes.Buffer
			var err = this.createFromElement(v.ElementItems, &temp, sqlArgMap)
			if err != nil {
				return err
			}
			sql.Write(temp.Bytes())
			loopChildItem = false
			break
		case Element_when:
			//if element
			var expression = v.Propertys[`test`]
			var result, err = this.doIfElement(expression, sqlArgMap, evaluateParameters)
			if err != nil {
				return err
			}
			if result {
				//test > true,write sql string
				var replaceSql, err = replaceArg(v.DataString, defaultArgMap, this.sqlArgTypeConvert, &this.expressionEngineProxy)
				if err != nil {
					return err
				}
				sql.WriteString(replaceSql)
				if loopChildItem && v.ElementItems != nil && len(v.ElementItems) > 0 {
					var err = this.createFromElement(v.ElementItems, sql, sqlArgMap)
					if err != nil {
						return err
					}
				}
				breakChildItem = true
			} else {
				// test > fail ,end loop
				loopChildItem = false
				break
			}
			break
		case Element_otherwise:
			if loopChildItem && v.ElementItems != nil && len(v.ElementItems) > 0 {
				var err = this.createFromElement(v.ElementItems, sql, sqlArgMap)
				if err != nil {
					return err
				}
			}
			breakChildItem = true
			break
		case Element_where:
			var err = this.elementTrim(&loopChildItem, v.ElementItems, sqlArgMap, DefaultWhereElement_Prefix, "", DefaultWhereElement_PrefixOverrides, "", sql)
			if err != nil {
				return err
			}
			break
		case Element_Include:
			//include have child elements,just break
			break
		default:
			panic("[GoMybatis] find not support element! " + v.ElementType)
		}
		if breakChildItem {
			break
		}
		if loopChildItem && v.ElementItems != nil && len(v.ElementItems) > 0 {
			var err = this.createFromElement(v.ElementItems, sql, sqlArgMap)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (this *GoMybatisSqlBuilder) doIfElement(expression string, param map[string]SqlArg, evaluateParameters map[string]interface{}) (bool, error) {
	//this.repleaceExpression(expression, param)
	ifElementevalExpression, err := this.expressionEngineProxy.Lexer(expression)
	if err != nil {
		return false, err
	}
	if evaluateParameters == nil {
		evaluateParameters = this.expressParamterMap(param, this.expressionTypeConvert)
	}
	result, err := this.expressionEngineProxy.Eval(ifElementevalExpression, evaluateParameters, 0)
	if err != nil {
		err = utils.NewError("SqlBuilder", "[GoMybatis] <test `", expression, `> fail,`, err.Error())
		return false, err
	}
	return result.(bool), nil
}

func (this *GoMybatisSqlBuilder) bindBindElementArg(args map[string]SqlArg, item ElementItem, typeConvert SqlArgTypeConvert) map[string]SqlArg {
	var name = item.Propertys["name"]
	var value = item.Propertys["value"]
	if name == "" {
		panic(`[GoMybatis] element <bind name = ""> name can not be nil!`)
	}
	if value == "" {
		args[name] = SqlArg{
			Value: fmt.Sprint(value),
			Type:  StringType,
		}
		return args
	}
	bindEvalExpression, err := this.expressionEngineProxy.Lexer(value)
	if err != nil {
		return args
	}
	var evaluateParameters = this.expressParamterMap(args, this.expressionTypeConvert)
	result, err := this.expressionEngineProxy.Eval(bindEvalExpression, evaluateParameters, 0)
	if err != nil {
		//TODO send log bind fail
		return args
	}
	args[name] = SqlArg{
		Value: fmt.Sprint(result),
		Type:  StringType,
	}
	return args
}

//scan params
func (this *GoMybatisSqlBuilder) expressParamterMap(parameters map[string]SqlArg, typeConvert ExpressionTypeConvert) map[string]interface{} {
	var newMap = make(map[string]interface{})
	for k, obj := range parameters {
		var value = obj.Value
		if typeConvert != nil {
			value = typeConvert.Convert(obj)
		}
		newMap[k] = value
	}
	return newMap
}

//scan params
func (this *GoMybatisSqlBuilder) sqlParamterMap(parameters map[string]SqlArg, typeConvert SqlArgTypeConvert) map[string]string {
	var newMap = make(map[string]string)
	for k, obj := range parameters {
		var value = obj.Value
		if typeConvert != nil {
			value = typeConvert.Convert(obj)
		}
		newMap[k] = fmt.Sprint(value)
	}
	return newMap
}



//trim处理element
func (this *GoMybatisSqlBuilder) elementTrim(loopChildItem *bool, items []ElementItem, param map[string]SqlArg, prefix string, suffix string, prefixOverrides string, suffixOverrides string, sql *bytes.Buffer) error {
	if *loopChildItem && items != nil && len(items) > 0 {
		var tempTrimSql bytes.Buffer
		var err = this.createFromElement(items, &tempTrimSql, param)
		if err != nil {
			return err
		}
		var tempTrimSqlString = strings.Trim(tempTrimSql.String(), " ")
		if prefixOverrides != "" {
			var prefixOverridesArray = strings.Split(prefixOverrides, "|")
			if len(prefixOverridesArray) > 0 {
				for _, v := range prefixOverridesArray {
					tempTrimSqlString = strings.TrimPrefix(tempTrimSqlString, v)
				}
			}
		}
		if suffixOverrides != "" {
			var suffixOverrideArray = strings.Split(suffixOverrides, "|")
			if len(suffixOverrideArray) > 0 {
				for _, v := range suffixOverrideArray {
					tempTrimSqlString = strings.TrimSuffix(tempTrimSqlString, v)
				}
			}
		}
		var newBuffer bytes.Buffer
		newBuffer.WriteString(` `)
		newBuffer.WriteString(prefix)
		newBuffer.WriteString(` `)
		newBuffer.WriteString(tempTrimSqlString)
		newBuffer.WriteString(` `)
		newBuffer.WriteString(suffix)
		sql.Write(newBuffer.Bytes())
		*loopChildItem = false
	}
	return nil
}

func (this *GoMybatisSqlBuilder) makeArgInterfaceMap(args map[string]SqlArg) map[string]interface{} {
	var m = make(map[string]interface{})
	if args != nil {
		for k, v := range args {
			m[k] = v.Value
		}
	}
	return m
}

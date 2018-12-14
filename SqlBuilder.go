package GoMybatis

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/zhuxiujia/GoMybatis/lib/github.com/Knetic/govaluate"
	"log"
	"reflect"
	"strings"
)

type SqlBuilder interface {
	BuildSql(paramMap map[string]SqlArg, mapperXml *MapperXml, enableLog bool) (string, error)
}

type GoMybatisSqlBuilder struct {
	SqlBuilder
	ExpressionTypeConvert ExpressionTypeConvert
	SqlArgTypeConvert     SqlArgTypeConvert
}

func (this GoMybatisSqlBuilder) New(ExpressionTypeConvert ExpressionTypeConvert, SqlArgTypeConvert SqlArgTypeConvert) GoMybatisSqlBuilder {
	this.ExpressionTypeConvert = ExpressionTypeConvert
	this.SqlArgTypeConvert = SqlArgTypeConvert
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

func (this GoMybatisSqlBuilder) createFromElement(itemTree []ElementItem, sql *bytes.Buffer, param map[string]SqlArg) error {
	if this.SqlArgTypeConvert == nil || this.ExpressionTypeConvert == nil {
		panic("[GoMybatis] GoMybatisSqlBuilder.SqlArgTypeConvert and GoMybatisSqlBuilder.ExpressionTypeConvert can not be nil!")
	}
	//test表达式参数map
	var evaluateParameters map[string]interface{}
	for _, v := range itemTree {
		var loopChildItem = true
		var breakChildItem = false
		switch v.ElementType {
		case Element_bind:
			//bind,param args change!need update
			param = this.bindBindElementArg(param, v, this.SqlArgTypeConvert)
			if evaluateParameters != nil {
				evaluateParameters = this.expressParamterMap(param, this.ExpressionTypeConvert)
			}
			break
		case Element_String:
			//string element
			sql.WriteString(replaceArg(v.DataString, param, this.SqlArgTypeConvert))
			break
		case Element_If:
			//if element
			var expression = v.Propertys[`test`]
			var result, err = this.doIfElement(&expression, param, evaluateParameters)
			if err != nil {
				return err
			}
			if result {
				//test > true,write sql string
				var reps = replaceArg(v.DataString, param, this.SqlArgTypeConvert)
				sql.WriteString(reps)
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
			var err = this.elementTrim(&loopChildItem, v.ElementItems, param, prefix, suffix, prefixOverrides, suffixOverrides, sql)
			if err != nil {
				return err
			}
			break
		case Element_Set:
			if loopChildItem && v.ElementItems != nil && len(v.ElementItems) > 0 {
				var trim bytes.Buffer
				var err = this.createFromElement(v.ElementItems, &trim, param)
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
			var datas = param[collection].Value
			var collectionValue = reflect.ValueOf(datas)
			var collectionValueLen = collectionValue.Len()
			if collectionValueLen == 0 {
				continue
			}
			for i := 0; i < collectionValueLen; i++ {
				var collectionItem = collectionValue.Index(i)
				var tempArgMap = make(map[string]SqlArg)
				for k, v := range param {
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
			var err = this.createFromElement(v.ElementItems, &temp, param)
			if err != nil {
				return err
			}
			sql.Write(temp.Bytes())
			loopChildItem = false
			break
		case Element_when:
			//if element
			var expression = v.Propertys[`test`]
			var result, err = this.doIfElement(&expression, param, evaluateParameters)
			if err != nil {
				return err
			}
			if result {
				//test > true,write sql string
				var reps = replaceArg(v.DataString, param, this.SqlArgTypeConvert)
				sql.WriteString(reps)
				if loopChildItem && v.ElementItems != nil && len(v.ElementItems) > 0 {
					var err = this.createFromElement(v.ElementItems, sql, param)
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
				var err = this.createFromElement(v.ElementItems, sql, param)
				if err != nil {
					return err
				}
			}
			breakChildItem = true
			break
		case Element_where:
			var prefix = "where"
			var suffix = ""
			var prefixOverrides = "and |or |And |Or "
			var suffixOverrides = ""
			var err = this.elementTrim(&loopChildItem, v.ElementItems, param, prefix, suffix, prefixOverrides, suffixOverrides, sql)
			if err != nil {
				return err
			}
			break
		default:
			panic("[GoMybatis] find not support element! " + v.ElementType)
		}
		if breakChildItem {
			break
		}
		if loopChildItem && v.ElementItems != nil && len(v.ElementItems) > 0 {
			var err = this.createFromElement(v.ElementItems, sql, param)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (this GoMybatisSqlBuilder) doIfElement(expression *string, param map[string]SqlArg, evaluateParameters map[string]interface{}) (bool, error) {
	this.repleaceExpression(expression, param)
	evalExpression, err := govaluate.NewEvaluableExpression(*expression)
	if err != nil {
		return false, err
	}
	if evaluateParameters == nil {
		evaluateParameters = this.expressParamterMap(param, this.ExpressionTypeConvert)
	}
	result, err := evalExpression.Evaluate(evaluateParameters)
	if err != nil {
		var buffer bytes.Buffer
		buffer.WriteString("[GoMybatis] <test `")
		buffer.WriteString(*expression)
		buffer.WriteString(`> fail,`)
		buffer.WriteString(err.Error())
		err = errors.New(buffer.String())
		return false, err
	}
	return result.(bool), nil
}

func (this GoMybatisSqlBuilder) bindBindElementArg(args map[string]SqlArg, item ElementItem, typeConvert SqlArgTypeConvert) map[string]SqlArg {
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
	evalExpression, err := govaluate.NewEvaluableExpression(value)
	if err != nil {
		return args
	}
	var evaluateParameters = this.expressParamterMap(args, this.ExpressionTypeConvert)
	result, err := evalExpression.Evaluate(evaluateParameters)
	if err != nil {
		return args
	}
	args[name] = SqlArg{
		Value: fmt.Sprint(result),
		Type:  StringType,
	}
	return args
}

//表达式 ''转换为 0
func (this GoMybatisSqlBuilder) expressionToIfZeroExpression(expression string, param map[string]SqlArg) string {
	for k, v := range param {
		if strings.Contains(expression, k) && v.Type.Kind() != reflect.String {
			expression = strings.Replace(expression, `''`, `0`, -1)
			return expression
		}
	}

	return expression
}

//scan params
func (this GoMybatisSqlBuilder) expressParamterMap(parameters map[string]SqlArg, typeConvert ExpressionTypeConvert) map[string]interface{} {
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
func (this GoMybatisSqlBuilder) sqlParamterMap(parameters map[string]SqlArg, typeConvert SqlArgTypeConvert) map[string]string {
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

func (this GoMybatisSqlBuilder) split(str *string) (stringItems []string) {
	if str == nil || *str == "" {
		return nil
	}
	var andStrings = strings.Split(*str, " && ")
	if andStrings == nil {
		return nil
	}
	var newStrings []string
	for _, v := range andStrings {
		var orStrings = strings.Split(v, " || ")
		if orStrings == nil {
			continue
		}
		for _, orStr := range orStrings {
			if newStrings == nil {
				newStrings = make([]string, 0)
			}
			if orStr == "" {
				continue
			}
			newStrings = append(newStrings, orStr)
		}
	}
	return newStrings
}

//替换表达式中的值 and,or,参数 替换为实际值
func (this GoMybatisSqlBuilder) repleaceExpression(expression *string, param map[string]SqlArg) {
	if expression == nil || *expression == "" {
		return
	}
	*expression = strings.Replace(*expression, ` and `, " && ", -1)
	*expression = strings.Replace(*expression, ` or `, " || ", -1)
	var newStrings = this.split(expression)

	for _, expressionItem := range newStrings {
		var NewExpression = this.expressionToIfZeroExpression(expressionItem, param)
		*expression = strings.Replace(*expression, expressionItem, NewExpression, -1)
	}
}

//trim处理element
func (this GoMybatisSqlBuilder) elementTrim(loopChildItem *bool, items []ElementItem, param map[string]SqlArg, prefix string, suffix string, prefixOverrides string, suffixOverrides string, sql *bytes.Buffer) error {
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

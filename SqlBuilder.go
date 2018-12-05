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
		if v.ElementType == Element_bind {
			//bind,param args change!need update
			param = this.bindBindElementArg(param, v, this.SqlArgTypeConvert)
			if evaluateParameters != nil {
				evaluateParameters = this.expressParamterMap(param, this.ExpressionTypeConvert)
			}
		} else if v.ElementType == Element_String {
			//string element
			sql.WriteString(replaceArg(v.DataString, param, this.SqlArgTypeConvert))
		} else if v.ElementType == Element_If {
			//if element
			var expression = v.Propertys[`test`]
			this.repleaceExpression(&expression, param)
			evalExpression, err := govaluate.NewEvaluableExpression(expression)
			if err != nil {
				return err
			}
			if evaluateParameters == nil {
				evaluateParameters = this.expressParamterMap(param, this.ExpressionTypeConvert)
			}
			result, err := evalExpression.Evaluate(evaluateParameters)
			if err != nil {
				var buffer bytes.Buffer
				buffer.WriteString("[GoMybatis] <test `")
				buffer.WriteString(expression)
				buffer.WriteString(`> fail,`)
				buffer.WriteString(err.Error())
				err = errors.New(buffer.String())
				return err
			}
			if result.(bool) {
				//test > true,write sql string
				var reps = replaceArg(v.DataString, param, this.SqlArgTypeConvert)
				sql.WriteString(reps)
			} else {
				// test > fail ,end loop
				loopChildItem = false
				break
			}
		} else if v.ElementType == Element_Trim {
			var prefix = v.Propertys[`prefix`]
			var suffix = v.Propertys[`suffix`]
			var suffixOverrides = v.Propertys[`suffixOverrides`]
			var prefixOverrides = v.Propertys[`prefixOverrides`]
			if loopChildItem && v.ElementItems != nil && len(v.ElementItems) > 0 {
				var tempTrimSql bytes.Buffer
				var err = this.createFromElement(v.ElementItems, &tempTrimSql, param)
				if err != nil {
					return err
				}
				var tempTrimSqlString = strings.Trim(strings.Trim(strings.Trim(tempTrimSql.String(), " "), suffixOverrides), prefixOverrides)
				var newBuffer bytes.Buffer
				newBuffer.WriteString(` `)
				newBuffer.WriteString(prefix)
				newBuffer.WriteString(` `)
				newBuffer.WriteString(tempTrimSqlString)
				newBuffer.WriteString(` `)
				newBuffer.WriteString(suffix)
				sql.Write(newBuffer.Bytes())
				loopChildItem = false
			}
		} else if v.ElementType == Element_Set {
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
		} else if v.ElementType == Element_Foreach {
			var collection = v.Propertys[`collection`]
			var index = v.Propertys[`index`]
			var item = v.Propertys[`item`]
			var open = v.Propertys[`open`]
			var close = v.Propertys[`close`]
			var separator = v.Propertys[`separator`]

			if item == "" {
				item = "item"
			}
			if index == "" {
				index = "index"
			}
			if collection == "" {
				panic(`[GoMybatis] collection value can not be "" in <foreach collection=""> !`)
			}

			var tempSql bytes.Buffer
			var datas = param[collection].Value
			var collectionValue = reflect.ValueOf(datas)
			var collectionValueLen = collectionValue.Len()
			if collectionValueLen > 0 {
				for i := 0; i < collectionValueLen; i++ {
					var collectionItem = collectionValue.Index(i)
					var tempArgMap = make(map[string]SqlArg)
					for k, v := range param {
						tempArgMap[k] = v
					}
					tempArgMap[item] = SqlArg{
						Value: collectionItem.Interface(),
						Type:  collectionItem.Type(),
					}
					tempArgMap[index] = SqlArg{
						Value: index,
						Type:  IntType,
					}
					if loopChildItem && v.ElementItems != nil && len(v.ElementItems) > 0 {
						var err = this.createFromElement(v.ElementItems, &tempSql, tempArgMap)
						if err != nil {
							return err
						}
					}
				}
			}
			var newTempSql bytes.Buffer
			newTempSql.WriteString(open)
			newTempSql.Write(tempSql.Bytes())
			newTempSql.WriteString(close)
			var tempSqlString = strings.Trim(strings.Trim(newTempSql.String(), " "), separator)
			tempSql.Reset()
			tempSql.WriteString(` `)
			tempSql.WriteString(tempSqlString)
			sql.Write(tempSql.Bytes())
			loopChildItem = false
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

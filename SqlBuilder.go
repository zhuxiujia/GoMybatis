package GoMybatis

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/zhuxiujia/GoMybatis/lib/github.com/Knetic/govaluate"
	"reflect"
	"strings"
)

type SqlBuilder interface {
	BuildSql(paramMap map[string]SqlArg, mapperXml MapperXml) (string, error)
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

func (this GoMybatisSqlBuilder) BuildSql(paramMap map[string]SqlArg, mapperXml MapperXml) (string, error) {
	var sql bytes.Buffer
	sql, err := this.createFromElement(mapperXml.ElementItems, sql, paramMap)
	if err != nil {
		return sql.String(), err
	}
	//log.Println("[GoMybatis] Preparing sql ==> ", sql.String())
	return sql.String(), nil
}

func (this GoMybatisSqlBuilder) createFromElement(itemTree []ElementItem, sql bytes.Buffer, param map[string]SqlArg) (result bytes.Buffer, err error) {
	if this.SqlArgTypeConvert == nil || this.ExpressionTypeConvert == nil {
		panic("[GoMybatis] GoMybatisSqlBuilder.SqlArgTypeConvert and GoMybatisSqlBuilder.ExpressionTypeConvert can not be nil!")
	}
	for _, v := range itemTree {
		var loopChildItem = true
		if v.ElementType == Element_String {
			//string element
			sql.WriteString(replaceArg(v.DataString, param, this.SqlArgTypeConvert))
		} else if v.ElementType == Element_If {
			//if element
			var test = v.Propertys[`test`]
			var andStrings = strings.Split(test, ` and `)
			for index, expression := range andStrings {
				//test表达式解析
				var evaluateParameters = this.scanParamterMap(param, this.ExpressionTypeConvert)
				expression = this.expressionToIfZeroExpression(expression, param)
				evalExpression, err := govaluate.NewEvaluableExpression(expression)
				if err != nil {
					fmt.Println(err)
				}
				result, err := evalExpression.Evaluate(evaluateParameters)
				if err != nil {
					var buffer bytes.Buffer
					buffer.WriteString("[GoMybatis] <test `")
					buffer.WriteString(expression)
					buffer.WriteString(`> fail,`)
					buffer.WriteString(err.Error())
					err = errors.New(buffer.String())
					return sql, err
				}
				if result.(bool) {
					//test表达式成立
					if index == (len(andStrings) - 1) {
						var reps = replaceArg(v.DataString, param, this.SqlArgTypeConvert)
						sql.WriteString(reps)
					}
				} else {
					loopChildItem = false
					break
				}
			}
		} else if v.ElementType == Element_Trim {
			var prefix = v.Propertys[`prefix`]
			var suffix = v.Propertys[`suffix`]
			var suffixOverrides = v.Propertys[`suffixOverrides`]
			var prefixOverrides = v.Propertys[`prefixOverrides`]
			if v.ElementItems != nil && len(v.ElementItems) > 0 && loopChildItem {
				var tempTrimSql bytes.Buffer
				tempTrimSql, err = this.createFromElement(v.ElementItems, tempTrimSql, param)
				if err != nil {
					return tempTrimSql, err
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
			if v.ElementItems != nil && len(v.ElementItems) > 0 && loopChildItem {
				var trim bytes.Buffer
				trim, err = this.createFromElement(v.ElementItems, trim, param)
				if err != nil {
					return trim, err
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
			var tempSql bytes.Buffer
			var datas = param[collection].Value
			var collectionValue = reflect.ValueOf(datas)
			if collectionValue.Len() > 0 {
				for i := 0; i < collectionValue.Len(); i++ {
					var dataItem = collectionValue.Index(i).Interface()
					var tempArgMap = make(map[string]SqlArg)
					tempArgMap[item] = SqlArg{
						Value: dataItem,
						Type:  collectionValue.Index(i).Type(),
					}
					tempArgMap[index] = SqlArg{
						Value: index,
						Type:  IntType,
					}
					for k, v := range param {
						tempArgMap[k] = v
					}
					if v.ElementItems != nil && len(v.ElementItems) > 0 && loopChildItem {
						tempSql, err = this.createFromElement(v.ElementItems, tempSql, tempArgMap)
						if err != nil {
							return tempSql, err
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
		if v.ElementItems != nil && len(v.ElementItems) > 0 && loopChildItem {
			sql, err = this.createFromElement(v.ElementItems, sql, param)
			if err != nil {
				return sql, err
			}
		}
	}
	return sql, nil
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
func (this GoMybatisSqlBuilder) scanParamterMap(parameters map[string]SqlArg, typeConvert ExpressionTypeConvert) map[string]interface{} {
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

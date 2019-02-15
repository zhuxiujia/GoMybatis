package GoMybatis

import (
	"bytes"
	"github.com/zhuxiujia/GoMybatis/utils"
	"reflect"
	"strings"
)

var equalOperator = []string{"/", "+", "-", "*", "**", "|", "^", "&", "%", "<", ">", ">=", "<=", " in ", " not in ", " or ", "||", " and ", "&&", "==", "!="}

type GoMybatisTempleteDecoder struct {
}

type LogicDeleteData struct {
	Enable         bool
	Deleted_value  string
	Undelete_value string
	Column         string
	Property       string
	LangType       string
}

func (it *GoMybatisTempleteDecoder) DecodeTree(tree map[string]*MapperXml, beanType reflect.Type) error {
	if tree == nil {
		return utils.NewError("GoMybatisTempleteDecoder", "decode data map[string]*MapperXml cant be nil!")
	}
	if beanType != nil {
		if beanType.Kind() == reflect.Ptr {
			beanType = beanType.Elem()
		}
	}
	for _, v := range tree {
		var method *reflect.StructField
		if beanType != nil {
			if isMethodElement(v.Tag) {
				var upperId = utils.UpperFieldFirstName(v.Id)
				m, haveMethod := beanType.FieldByName(upperId)
				if haveMethod {
					method = &m
				}
			}
		}
		it.Decode(method, v, tree)
	}
	return nil
}

func (it *GoMybatisTempleteDecoder) Decode(method *reflect.StructField, mapper *MapperXml, tree map[string]*MapperXml) error {

	switch mapper.Tag {

	case "selectTemplete":
		mapper.Tag = Element_Select

		var tables = mapper.Propertys["tables"]
		var columns = mapper.Propertys["columns"]
		var wheres = mapper.Propertys["wheres"]

		var resultMap = mapper.Propertys["resultMap"]
		if resultMap == "" {
			resultMap = "BaseResultMap"
		}
		var resultMapData = tree[resultMap]
		if resultMapData == nil {
			panic(utils.NewError("GoMybatisTempleteDecoder", "resultMap not define! id = ", resultMap))
		}
		var logic = it.decodeLogicDelete(resultMapData)

		var sql bytes.Buffer
		sql.WriteString("select ")
		if columns == "" {
			columns = "*"
		}
		sql.WriteString(columns)
		sql.WriteString(" from ")
		sql.WriteString(tables)
		if len(wheres) > 0 {
			sql.WriteString(" where ")
			mapper.ElementItems = append(mapper.ElementItems, ElementItem{
				ElementType: Element_String,
				DataString:  sql.String(),
			})
			//TODO decode wheres
			it.DecodeWheres(wheres, mapper, logic)
		}
		break

	case "insertTemplete":
		mapper.Tag = Element_Insert

		var tables = mapper.Propertys["tables"]
		var inserts = mapper.Propertys["inserts"]

		var resultMap = mapper.Propertys["resultMap"]
		if resultMap == "" {
			resultMap = "BaseResultMap"
		}
		if inserts == "" {
			inserts = "*?*"
		}

		var resultMapData = tree[resultMap]
		if resultMapData == nil {
			panic(utils.NewError("GoMybatisTempleteDecoder", "resultMap not define! id = ", resultMap))
		}
		var logic = it.decodeLogicDelete(resultMapData)

		var sql bytes.Buffer
		sql.WriteString("insert into ")
		sql.WriteString(tables)

		mapper.ElementItems = append(mapper.ElementItems, ElementItem{
			ElementType: Element_String,
			DataString:  sql.String(),
		})

		//add insert column
		var trimColumn = ElementItem{
			ElementType:  Element_Trim,
			Propertys:    map[string]string{"prefix": "(", "suffix": ")", "suffixOverrides": ","},
			ElementItems: []ElementItem{},
		}

		for _, v := range resultMapData.ElementItems {
			if logic.Enable && v.Propertys["property"] == logic.Property {
				trimColumn.ElementItems = append(trimColumn.ElementItems, ElementItem{
					ElementType: Element_String,
					DataString:  logic.Undelete_value + ",",
				})
				continue
			}
			if inserts == "*?*" {
				trimColumn.ElementItems = append(trimColumn.ElementItems, ElementItem{
					ElementType: Element_If,
					Propertys:   map[string]string{"test": it.makeIfNotNull(v.Propertys["property"])},
					ElementItems: []ElementItem{
						{
							ElementType: Element_String,
							DataString:  v.Propertys["column"] + ",",
						},
					},
				})
			} else if inserts == "*" {
				trimColumn.ElementItems = append(trimColumn.ElementItems, ElementItem{
					ElementType: Element_String,
					DataString:  v.Propertys["column"] + ",",
				})
			}
		}

		mapper.ElementItems = append(mapper.ElementItems, trimColumn)

		//add insert arg
		var trimArg = ElementItem{
			ElementType:  Element_Trim,
			Propertys:    map[string]string{"prefix": "values (", "suffix": ")", "suffixOverrides": ","},
			ElementItems: []ElementItem{},
		}
		for _, v := range resultMapData.ElementItems {
			if logic.Enable && v.Propertys["property"] == logic.Property {
				trimArg.ElementItems = append(trimArg.ElementItems, ElementItem{
					ElementType: Element_String,
					DataString:  logic.Undelete_value + ",",
				})
				continue
			}
			if inserts == "*?*" {
				trimArg.ElementItems = append(trimArg.ElementItems, ElementItem{
					ElementType: Element_If,
					Propertys:   map[string]string{"test": it.makeIfNotNull(v.Propertys["property"])},
					DataString:  "#{" + v.Propertys["property"] + "},",
				})
			} else if inserts == "*" {
				trimArg.ElementItems = append(trimArg.ElementItems, ElementItem{
					ElementType: Element_String,
					DataString:  "#{" + v.Propertys["property"] + "},",
				})
			}
		}
		mapper.ElementItems = append(mapper.ElementItems, trimArg)
		break
	case "updateTemplete":
		mapper.Tag = Element_Update

		var tables = mapper.Propertys["tables"]
		var columns = mapper.Propertys["sets"]
		var wheres = mapper.Propertys["wheres"]

		var resultMap = mapper.Propertys["resultMap"]
		if resultMap == "" {
			resultMap = "BaseResultMap"
		}

		var resultMapData = tree[resultMap]
		if resultMapData == nil {
			panic(utils.NewError("GoMybatisTempleteDecoder", "resultMap not define! id = ", resultMap))
		}

		var logic = it.decodeLogicDelete(resultMapData)

		var sql bytes.Buffer
		sql.WriteString("update set ")
		if columns == "" {
			mapper.ElementItems = append(mapper.ElementItems, ElementItem{
				ElementType: Element_String,
				DataString:  sql.String(),
			})
			sql.Reset()
			for _, v := range resultMapData.ElementItems {
				var column = v.Propertys["column"]
				var property = v.Propertys["property"]
				mapper.ElementItems = append(mapper.ElementItems, ElementItem{
					ElementType: Element_If,
					Propertys:   map[string]string{"test": it.makeIfNotNull(property)},
					DataString:  column + " = #{" + v.Propertys["property"] + "},",
				})
			}
		} else {
			mapper.ElementItems = append(mapper.ElementItems, ElementItem{
				ElementType: Element_String,
				DataString:  sql.String(),
			})
			sql.Reset()
			it.DecodeSets(columns, mapper, LogicDeleteData{})
		}
		sql.WriteString(" from ")
		sql.WriteString(tables)

		if len(wheres) > 0 || logic.Enable {
			sql.WriteString(" where ")
			mapper.ElementItems = append(mapper.ElementItems, ElementItem{
				ElementType: Element_String,
				DataString:  sql.String(),
			})
			it.DecodeWheres(wheres, mapper, logic)
		}
		break

	case "deleteTemplete":
		mapper.Tag = Element_Delete

		var tables = mapper.Propertys["tables"]
		var wheres = mapper.Propertys["wheres"]

		var resultMap = mapper.Propertys["resultMap"]
		if resultMap == "" {
			resultMap = "BaseResultMap"
		}

		var resultMapData = tree[resultMap]
		if resultMapData == nil {
			panic(utils.NewError("GoMybatisTempleteDecoder", "resultMap not define! id = ", resultMap))
		}

		var logic = it.decodeLogicDelete(resultMapData)
		if logic.Enable {
			//enable logic delete
			var sql bytes.Buffer
			sql.WriteString("update set ")

			mapper.ElementItems = append(mapper.ElementItems, ElementItem{
				ElementType: Element_String,
				DataString:  sql.String(),
			})
			sql.Reset()
			it.DecodeSets("", mapper, logic)

			sql.WriteString(" from ")
			sql.WriteString(tables)
			if len(wheres) > 0 {
				sql.WriteString(" where ")
				mapper.ElementItems = append(mapper.ElementItems, ElementItem{
					ElementType: Element_String,
					DataString:  sql.String(),
				})
				//TODO decode wheres
				it.DecodeWheres(wheres, mapper, logic)
			}
			break
		} else {
			//default delete  DELETE FROM `test`.`biz_activity` WHERE `id`='165';
			var sql bytes.Buffer
			sql.WriteString("delete from ")
			sql.WriteString(tables)
			if len(wheres) > 0 {
				sql.WriteString(" where ")
				mapper.ElementItems = append(mapper.ElementItems, ElementItem{
					ElementType: Element_String,
					DataString:  sql.String(),
				})
				//TODO decode wheres
				it.DecodeWheres(wheres, mapper, LogicDeleteData{})
			}
		}
	}

	return nil
}

//解码逗号分隔的where
func (it *GoMybatisTempleteDecoder) DecodeWheres(arg string, mapper *MapperXml, logic LogicDeleteData) {
	var wheres = strings.Split(arg, ",")
	for index, v := range wheres {
		var expressions = strings.Split(v, "?")
		if len(expressions) > 1 {
			//TODO have ?
			var newWheres bytes.Buffer
			if index > 0 {
				newWheres.WriteString(" and ")
			}
			newWheres.WriteString(expressions[1])
			var item = ElementItem{
				ElementType: Element_If,
				Propertys:   map[string]string{"test": it.makeIfNotNull(expressions[0])},
				DataString:  newWheres.String(),
			}
			mapper.ElementItems = append(mapper.ElementItems, item)
		} else {
			var newWheres bytes.Buffer
			if index > 0 {
				newWheres.WriteString(" and ")
			}
			newWheres.WriteString(v)
			var item = ElementItem{
				ElementType: Element_String,
				DataString:  newWheres.String(),
			}
			mapper.ElementItems = append(mapper.ElementItems, item)
		}
	}
	if logic.Enable == true {
		var appendAdd = ""
		if len(wheres) >= 1 {
			appendAdd = " and "
		}
		var item = ElementItem{
			ElementType: Element_String,
			DataString:  appendAdd + logic.Column + " = " + logic.Undelete_value,
		}
		mapper.ElementItems = append(mapper.ElementItems, item)
	}
}

func (it *GoMybatisTempleteDecoder) DecodeSets(arg string, mapper *MapperXml, logic LogicDeleteData) {
	var sets = strings.Split(arg, ",")
	for index, v := range sets {
		var expressions = strings.Split(v, "?")
		if len(expressions) > 1 {
			//TODO have ?
			var newWheres bytes.Buffer
			if index > 0 {
				newWheres.WriteString(",")
			}
			newWheres.WriteString(expressions[1])
			var item = ElementItem{
				ElementType: Element_If,
				Propertys:   map[string]string{"test": it.makeIfNotNull(expressions[0])},
				DataString:  newWheres.String(),
			}
			mapper.ElementItems = append(mapper.ElementItems, item)
		} else {
			var newWheres bytes.Buffer
			if index > 0 {
				newWheres.WriteString(" and ")
			}
			newWheres.WriteString(v)
			var item = ElementItem{
				ElementType: Element_String,
				DataString:  newWheres.String(),
			}
			mapper.ElementItems = append(mapper.ElementItems, item)
		}
	}
	if logic.Enable == true {
		var appendAdd = ""
		if len(sets) >= 1 {
			appendAdd = ","
		}
		var item = ElementItem{
			ElementType: Element_String,
			DataString:  appendAdd + logic.Column + " = " + logic.Deleted_value,
		}
		mapper.ElementItems = append(mapper.ElementItems, item)
	}
}

func (it *GoMybatisTempleteDecoder) makeIfNotNull(arg string) string {
	for _, v := range equalOperator {
		if strings.Contains(arg, v) {
			return arg
		}
	}
	return arg + ` != null`
}

func (it *GoMybatisTempleteDecoder) decodeLogicDelete(xml *MapperXml) LogicDeleteData {
	if xml == nil || len(xml.ElementItems) == 0 {
		return LogicDeleteData{}
	}
	var logicData = LogicDeleteData{}
	for _, v := range xml.ElementItems {
		if v.Propertys["logic_enable"] == "true" {
			logicData.Enable = true
			logicData.Deleted_value = v.Propertys["logic_deleted_value"]
			logicData.Undelete_value = v.Propertys["logic_undelete_value"]
			logicData.Column = v.Propertys["column"]
			logicData.Property = v.Propertys["property"]
			logicData.LangType = v.Propertys["langType"]
			break
		}
	}
	return logicData
}

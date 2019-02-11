package GoMybatis

import (
	"bytes"
	"strings"
)

var equalOperator = []string{"/", "+", "-", "*", "**", "|", "^", "&", "%", "<", ">", ">=", "<=", " in ", " not in ", " or ", "||", " and ", "&&", "==", "!="}

type GoMybatisTempleteDecoder struct {
}

func (it *GoMybatisTempleteDecoder) DecodeTree(tree map[string]*MapperXml) error {
	for _, v := range tree {
		it.Decode(v)
	}
	return nil
}

func (it *GoMybatisTempleteDecoder) Decode(mapper *MapperXml) error {

	var tables = mapper.Propertys["tables"]
	var columns = mapper.Propertys["columns"]
	var wheres = mapper.Propertys["wheres"]
	var sql bytes.Buffer
	if mapper.Tag == "selectTemplete" {
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
			sql.Reset()
			it.DecodeWheres(wheres, mapper)
		}

		if mapper.Id == "" {
			mapper.Id = mapper.Tag
		}
		mapper.Tag = Element_Select
	}

	return nil
}

//解码逗号分隔的where
func (it *GoMybatisTempleteDecoder) DecodeWheres(arg string, mapper *MapperXml) {
	var wheres = strings.Split(arg, ",")
	for index, v := range wheres {
		var expressions = strings.Split(v, "?")
		if len(expressions) > 1 {
			//TODO have ?
			var newWheres bytes.Buffer
			for k, v := range expressions {
				if k > 0 {
					if index > 0 {
						newWheres.WriteString(" and ")
					}
					newWheres.WriteString(v)
				}
			}
			var item = ElementItem{
				ElementType: Element_If,
				Propertys:   map[string]string{"test": it.convertEqualAction(expressions[0])},
				DataString:  newWheres.String(),
			}
			mapper.ElementItems = append(mapper.ElementItems, item)
		}
	}
}

func (it *GoMybatisTempleteDecoder) convertEqualAction(arg string) string {
	for _, v := range equalOperator {
		if strings.Contains(arg, v) {
			return arg
		}
	}
	return arg + ` != null`
}

package GoMybatis

import (
	"bytes"
	"strings"
)

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
				Propertys:   map[string]string{"test": expressions[0] + ` != null`},
				DataString:  newWheres.String(),
			}
			mapper.ElementItems = append(mapper.ElementItems, item)
		}
	}
}

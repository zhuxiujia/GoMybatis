package GoMybatis

import (
	"bytes"
	"fmt"
)

type GoMybatisTempleteDecoder struct {
}

func (it *GoMybatisTempleteDecoder) Decode(mapper *MapperXml) error {

	var table = mapper.Propertys["table"]
	var columns = mapper.Propertys["columns"]
	var wheres = mapper.Propertys["wheres"]
	fmt.Println(table)
	fmt.Println(columns)
	fmt.Println(wheres)

	//TODO decode table
	//TODO decode columns
	//TODO decode wheres

	var sql bytes.Buffer
	if mapper.Tag == "selectTemplete" {
		sql.WriteString("select ")
		sql.WriteString(columns)
		sql.WriteString(" from ")
		sql.WriteString(table)
		sql.WriteString(" where ")
		sql.WriteString(wheres)
		sql.WriteString(" ")
		mapper.ElementItems= append(mapper.ElementItems, ElementItem{
			DataString:sql.String(),
		})
	}

	return nil
}

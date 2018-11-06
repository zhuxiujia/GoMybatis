package GoMybatis

import "encoding/xml"

type ResultMap struct {
	XMLName xml.Name         `xml:"resultMap"`
	Id      IdProperty       `xml:"id"`
	Results []ResultProperty `xml:"result"`
}

type IdProperty struct {
	XMLName  xml.Name `xml:"id"`
	Column   string   `xml:"column,attr"`
	Property string   `xml:"property,attr"`
	JdbcType string   `xml:"jdbcType,attr"`
}
type ResultProperty struct {
	XMLName  xml.Name `xml:"result"`
	Column   string   `xml:"column,attr"`
	Property string   `xml:"property,attr"`
	JdbcType string   `xml:"jdbcType,attr"`
}

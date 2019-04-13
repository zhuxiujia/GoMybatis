package GoMybatis

import (
	"github.com/zhuxiujia/GoMybatis/lib/github.com/beevik/etree"
	"github.com/zhuxiujia/GoMybatis/utils"
	"reflect"
)

const Element_Mapper = "mapper"
const ID = `id`

func LoadMapperXml(bytes []byte) (items map[string]etree.Token) {
	utils.FixTestExpressionSymbol(&bytes)
	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(bytes); err != nil {
		panic(err)
	}
	items = make(map[string]etree.Token)
	root := doc.SelectElement(Element_Mapper)
	for _, s := range root.ChildElements() {
		if s.Tag == Element_Insert ||
			s.Tag == Element_Delete ||
			s.Tag == Element_Update ||
			s.Tag == Element_Select ||
			s.Tag == Element_ResultMap ||
			s.Tag == Element_Sql ||
			s.Tag == Element_Insert_Templete ||
			s.Tag == Element_Delete_Templete ||
			s.Tag == Element_Update_Templete ||
			s.Tag == Element_Select_Templete {
			var elementID = s.SelectAttrValue(ID, "")

			if elementID == "" {
				//如果id不存在，id设置为tag
				elementID = s.Tag
			}
			if elementID != "" {
				var oldItem = items[elementID]
				if oldItem != nil {
					panic("[GoMybatis] element Id can not repeat in xml! elementId=" + elementID)
				}
			}
			items[elementID] = s
		}
	}
	for _, mapperXml := range items {
		var typeString = reflect.TypeOf(mapperXml).String()
		if typeString == "*etree.Element" {
			var el = mapperXml.(*etree.Element)
			for _, v := range el.ChildElements() {
				includeElementReplace(v, &items)
			}
		}
	}
	return items
}

func includeElementReplace(xml *etree.Element, xmlMap *map[string]etree.Token) {
	if xml.Tag == Element_Include {
		var refid = xml.SelectAttr("refid").Value
		if refid == "" {
			panic(`[GoMybatis] xml <includ refid=""> 'refid' can not be ""`)
		}
		var mapperXml = (*xmlMap)[refid]
		if mapperXml == nil {
			panic(`[GoMybatis] xml <includ refid="` + refid + `"> element can not find !`)
		}
		if xml != nil {
			(*xml).Child = mapperXml.(*etree.Element).Child
		}
	}
	if xml.Child != nil {
		for _, v := range xml.ChildElements() {
			includeElementReplace(v, xmlMap)
		}
	}
}

////标签上下级关系检查
//func elementRuleCheck(fatherElement *etree.Element, childElementItem ElementItem) {
//	if fatherElement.Tag != Element_choose && (childElementItem.ElementType == Element_when || childElementItem.ElementType == Element_otherwise) {
//		panic("[GoMybatis] find element <" + childElementItem.ElementType + "> not in <choose>!")
//	}
//}

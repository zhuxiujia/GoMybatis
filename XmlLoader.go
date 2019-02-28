package GoMybatis

import (
	"github.com/zhuxiujia/GoMybatis/lib/github.com/beevik/etree"
	"github.com/zhuxiujia/GoMybatis/utils"
)

const Element_Mapper = "mapper"
const ID = `id`

func LoadMapperXml(bytes []byte) (items map[string]*etree.Element) {
	utils.FixTestExpressionSymbol(&bytes)
	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(bytes); err != nil {
		panic(err)
	}
	items = make(map[string]*etree.Element)
	root := doc.SelectElement(Element_Mapper)
	for _, s := range root.ChildElements() {
		var attrMap = attrToProperty(s.Attr)
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
			var elementID = attrMap[ID]

			if elementID == "" {
				//如果id不存在，id设置为tag
				attrMap[ID] = s.Tag
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
	for itemsIndex, mapperXml := range items {
		for key, v := range mapperXml.ChildElements() {
			var isChanged = includeElementReplace(v, &items)
			if isChanged {
				mapperXml.Child[key] = v
			}
		}
		items[itemsIndex] = mapperXml
	}
	return items
}

func includeElementReplace(xml *etree.Element, xmlMap *map[string]*etree.Element) bool {
	var changed = false
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
			(*xml).Child = mapperXml.Child
			changed = true
		}
	}
	if xml.Child != nil {
		for index, v := range xml.ChildElements() {
			var isChanged = includeElementReplace(v, xmlMap)
			if isChanged {
				xml.Child[index] = v
			}
		}
	}
	return changed
}

func attrToProperty(attrs []etree.Attr) map[string]string {
	var m = make(map[string]string)
	for _, v := range attrs {
		m[v.Key] = v.Value
	}
	return m
}

////标签上下级关系检查
//func elementRuleCheck(fatherElement *etree.Element, childElementItem ElementItem) {
//	if fatherElement.Tag != Element_choose && (childElementItem.ElementType == Element_when || childElementItem.ElementType == Element_otherwise) {
//		panic("[GoMybatis] find element <" + childElementItem.ElementType + "> not in <choose>!")
//	}
//}

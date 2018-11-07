package GoMybatis

import (
	"github.com/zhuxiujia/GoMybatis/lib/github.com/beevik/etree"
	"reflect"
	"strings"
)

type MapperXml struct {
	Tag          string
	Id           string
	ElementItems []ElementItem
}

type ElementItem struct {
	ElementType  string
	Propertys    map[string]string
	DataString   string
	ElementItems []ElementItem
}
//读取xml
func LoadMapperXml(bytes []byte) (items []MapperXml) {
	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(bytes); err != nil {
		panic(err)
	}
	root := doc.SelectElement("mapper")
	for _, s := range root.ChildElements() {
		var attrMap = attrToProperty(s.Attr)
		var elItems = loop(s)
		if s.Tag == Insert || s.Tag == Delete || s.Tag == Update || s.Tag == Select {
			items = append(items, MapperXml{
				Tag:          s.Tag,
				Id:           attrMap[`id`],
				ElementItems: elItems,
			})
		}
	}
	return items
}

func attrToProperty(attrs []etree.Attr) map[string]string {
	var m = make(map[string]string)
	for _, v := range attrs {
		m[v.Key] = v.Value
	}
	return m
}

func loop(element *etree.Element) []ElementItem {
	var els = make([]ElementItem, 0)
	for _, el := range element.Child {
		var typeString = reflect.ValueOf(el).Type().String()
		if typeString == `*etree.CharData` {
			var d = el.(*etree.CharData)
			var str = d.Data
			if str == "" {
				continue
			}
			str = strings.Replace(str, "\n", "", -1)
			str = strings.Replace(str, "\t", "", -1)
			str = strings.Trim(str, " ")
			if str != "" {
				str = str + " "
				var elementItem = ElementItem{
					ElementType: String,
					DataString:  str,
				}
				els = append(els, elementItem)
			}
		} else if typeString == `*etree.Element` {
			var e = el.(*etree.Element)
			var element = ElementItem{
				ElementType:  e.Tag,
				ElementItems: make([]ElementItem, 0),
				Propertys:    attrToProperty(e.Attr),
			}
			if len(e.Child) > 0 {
				var loopEls = loop(e)
				for _, item := range loopEls {
					element.ElementItems = append(element.ElementItems, item)
				}
			}
			els = append(els, element)
		}
	}
	return els
}
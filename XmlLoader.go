package GoMybatis

import (
	"bytes"
	"github.com/zhuxiujia/GoMybatis/lib/github.com/beevik/etree"
	"reflect"
	"strings"
)

const EtreeCharData = `*etree.CharData`
const EtreeElement = `*etree.Element`

const Element_Mapper = "mapper"
const ID = `id`

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
	root := doc.SelectElement(Element_Mapper)
	for _, s := range root.ChildElements() {
		var attrMap = attrToProperty(s.Attr)
		var elItems = loop(s)
		if s.Tag == Element_Insert ||
			s.Tag == Element_Delete ||
			s.Tag == Element_Update ||
			s.Tag == Element_Select ||
			s.Tag == Element_ResultMap{
			items = append(items, MapperXml{
				Tag:          s.Tag,
				Id:           attrMap[ID],
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
		if typeString == EtreeCharData {
			var d = el.(*etree.CharData)
			var str = d.Data
			if str == "" {
				continue
			}
			str = strings.Replace(str, "\n", "", -1)
			str = strings.Replace(str, "\t", "", -1)
			str = strings.Trim(str, " ")
			if str != "" {
				var buf bytes.Buffer
				buf.WriteString(" ")
				buf.WriteString(str)
				var elementItem = ElementItem{
					ElementType: Element_String,
					DataString:  buf.String(),
				}
				els = append(els, elementItem)
			}
		} else if typeString == EtreeElement {
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

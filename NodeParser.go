package GoMybatis

import (
	"github.com/zhuxiujia/GoMybatis/lib/github.com/beevik/etree"
	"reflect"
)

const EtreeCharData = `*etree.CharData`
const EtreeElement = `*etree.Element`

//节点解析器
type NodeParser struct {
	holder NodeConfigHolder
}

//界面为node
func (it NodeParser) ParserNodes(mapperXml []etree.Token) []Node {
	if it.holder.proxy == nil {
		panic("NodeParser need a *ExpressionEngineProxy{}!")
	}
	var nodes = []Node{}
	for _, item := range mapperXml {
		var node Node
		var typeString = reflect.TypeOf(item).String()

		if typeString == EtreeCharData {
			charData := item.(*etree.CharData)
			n := NodeString{
				value:               charData.Data,
				t:                   NString,
				expressMap:          FindAllExpressConvertString(charData.Data), //表达式需要替换的string
				noConvertExpressMap: FindAllExpressString(charData.Data),
				holder:              &it.holder,
			}
			if len(n.expressMap) == 0 {
				n.expressMap = nil
			}
			node = &n
			continue
		}

		var v = item.(*etree.Element)
		var childItems = v.Child
		//start
		switch v.Tag {
		//case "string":
		//	n := NodeString{
		//		value:               v.Text(),
		//		t:                   NString,
		//		expressMap:          FindAllExpressConvertString(v.Text()), //表达式需要替换的string
		//		noConvertExpressMap: FindAllExpressString(v.Text()),
		//		holder:              &it.holder,
		//	}
		//	if len(n.expressMap) == 0 {
		//		n.expressMap = nil
		//	}
		//	node = &n
		//	break
		case "if":
			n := NodeIf{
				t:      NIf,
				test:   v.SelectAttrValue("test", ""),
				childs: []Node{},
				holder: &it.holder,
			}
			if childItems != nil {
				var childNodes = it.ParserNodes(childItems)
				n.childs = append(n.childs, childNodes...)
			}
			node = &n
			break
		case "trim":
			n := NodeTrim{
				t:               NTrim,
				prefix:          []byte(v.SelectAttrValue("prefix", "")),
				suffix:          []byte(v.SelectAttrValue("suffix", "")),
				prefixOverrides: []byte(v.SelectAttrValue("prefixOverrides", "")),
				suffixOverrides: []byte(v.SelectAttrValue("suffixOverrides", "")),
				childs:          []Node{},
			}
			if childItems != nil {
				var childNodes = it.ParserNodes(childItems)
				n.childs = append(n.childs, childNodes...)
			}
			node = &n
			break
		case "set":
			n := NodeTrim{
				t:      NTrim,
				childs: []Node{},

				prefix:          []byte(" set "),
				suffix:          nil,
				prefixOverrides: []byte(","),
				suffixOverrides: []byte(","),
			}
			if childItems != nil {
				var childNodes = it.ParserNodes(childItems)
				n.childs = append(n.childs, childNodes...)
			}
			node = &n
			break
		case "foreach":
			n := NodeForEach{
				t:          NForEach,
				childs:     []Node{},
				collection: v.SelectAttrValue("collection", ""),
				index:      v.SelectAttrValue("index", ""),
				item:       v.SelectAttrValue("item", ""),
				open:       v.SelectAttrValue("open", ""),
				close:      v.SelectAttrValue("close", ""),
				separator:  v.SelectAttrValue("separator", ""),
			}
			if childItems != nil {
				var childNodes = it.ParserNodes(childItems)
				n.childs = append(n.childs, childNodes...)
			}
			node = &n
			break
		case "choose":
			n := NodeChoose{
				t:         NChoose,
				whenNodes: []Node{},
			}
			if childItems != nil {
				var childNodes = it.ParserNodes(childItems)
				for _, v := range childNodes {
					if v.Type() == NWhen {
						n.whenNodes = append(n.whenNodes, childNodes...)
					} else if v.Type() == NOtherwise {
						if n.otherwiseNode != nil {
							panic("element only support one Otherwise node!")
						}
						n.otherwiseNode = v
					} else {
						panic("not support element type:" + v.Type().ToString())
					}
				}

			} else {
				n.whenNodes = nil
				n.otherwiseNode = nil
			}
			node = &n
			break
		case "when":
			n := NodeWhen{
				t:      NOtherwise,
				childs: []Node{},
				test:   v.SelectAttrValue("test", ""),
				holder: &it.holder,
			}
			if childItems != nil {
				var childNodes = it.ParserNodes(childItems)
				n.childs = append(n.childs, childNodes...)
			}
			node = &n
			break
		case "otherwise":
			n := NodeOtherwise{
				t:      NOtherwise,
				childs: []Node{},
			}
			if childItems != nil {
				var childNodes = it.ParserNodes(childItems)
				n.childs = append(n.childs, childNodes...)
			}
			node = &n
			break
		case "where":
			n := NodeTrim{
				t:               NTrim,
				prefix:          []byte(DefaultWhereElement_Prefix),
				suffix:          []byte(v.SelectAttrValue("suffix", "")),
				prefixOverrides: []byte(DefaultWhereElement_PrefixOverrides),
				suffixOverrides: []byte(v.SelectAttrValue("suffixOverrides", "")),
				childs:          []Node{},
			}
			if childItems != nil {
				var childNodes = it.ParserNodes(childItems)
				n.childs = append(n.childs, childNodes...)
			}
			node = &n
			break
		case "bind":
			n := NodeBind{
				t:      NBind,
				value:  v.SelectAttrValue("value", ""),
				name:   v.SelectAttrValue("name", ""),
				holder: &it.holder,
			}
			node = &n
		}
		nodes = append(nodes, node)
	}
	return nodes
}

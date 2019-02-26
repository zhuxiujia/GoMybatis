package GoMybatis

//节点解析器
type NodeParser struct {
	holder NodeConfigHolder
}

//界面为node
func (it NodeParser) ParserNodes(mapperXml []ElementItem) []Node {
	if it.holder.proxy == nil {
		panic("NodeParser need a *ExpressionEngineProxy{}!")
	}
	var nodes = []Node{}
	for _, v := range mapperXml {
		var node Node

		switch v.ElementType {
		case "string":
			n := NodeString{
				value:               v.DataString,
				t:                   NString,
				expressMap:          FindAllExpressConvertString(v.DataString), //表达式需要替换的string
				noConvertExpressMap: FindAllExpressString(v.DataString),
				holder:              &it.holder,
			}
			if len(n.expressMap) == 0 {
				n.expressMap = nil
			}
			node = &n
			break
		case "if":
			n := NodeIf{
				t:      NIf,
				test:   v.Propertys["test"],
				childs: []Node{},
				holder: &it.holder,
			}
			if v.ElementItems != nil && len(v.ElementItems) > 0 {
				var childNodes = it.ParserNodes(v.ElementItems)
				n.childs = append(n.childs, childNodes...)
			} else {
				n.childs = nil
			}
			node = &n
			break
		case "trim":
			n := NodeTrim{
				t:               NTrim,
				prefix:          []byte(v.Propertys["prefix"]),
				suffix:          []byte(v.Propertys["suffix"]),
				prefixOverrides: []byte(v.Propertys["prefixOverrides"]),
				suffixOverrides: []byte(v.Propertys["suffixOverrides"]),
				childs:          []Node{},
			}
			if v.ElementItems != nil && len(v.ElementItems) > 0 {
				var childNodes = it.ParserNodes(v.ElementItems)
				n.childs = append(n.childs, childNodes...)
			} else {
				n.childs = nil
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
			if v.ElementItems != nil && len(v.ElementItems) > 0 {
				var childNodes = it.ParserNodes(v.ElementItems)
				n.childs = append(n.childs, childNodes...)
			} else {
				n.childs = nil
			}
			node = &n
			break
		case "foreach":
			n := NodeForEach{
				t:          NForEach,
				childs:     []Node{},
				collection: v.Propertys["collection"],
				index:      v.Propertys["index"],
				item:       v.Propertys["item"],
				open:       v.Propertys["open"],
				close:      v.Propertys["close"],
				separator:  v.Propertys["separator"],
			}
			if v.ElementItems != nil && len(v.ElementItems) > 0 {
				var childNodes = it.ParserNodes(v.ElementItems)
				n.childs = append(n.childs, childNodes...)
			} else {
				n.childs = nil
			}
			node = &n
			break
		case "choose":
			n := NodeChoose{
				t:         NChoose,
				whenNodes: []Node{},
			}
			if v.ElementItems != nil && len(v.ElementItems) > 0 {
				var childNodes = it.ParserNodes(v.ElementItems)
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
				test:   v.Propertys["test"],
				holder: &it.holder,
			}
			if v.ElementItems != nil && len(v.ElementItems) > 0 {
				var childNodes = it.ParserNodes(v.ElementItems)
				n.childs = append(n.childs, childNodes...)
			} else {
				n.childs = nil
			}
			node = &n
			break
		case "otherwise":
			n := NodeOtherwise{
				t:      NOtherwise,
				childs: []Node{},
			}
			if v.ElementItems != nil && len(v.ElementItems) > 0 {
				var childNodes = it.ParserNodes(v.ElementItems)
				n.childs = append(n.childs, childNodes...)
			} else {
				n.childs = nil
			}
			node = &n
			break
		case "where":
			n := NodeTrim{
				t:               NTrim,
				prefix:          []byte(DefaultWhereElement_Prefix),
				suffix:          []byte(v.Propertys["suffix"]),
				prefixOverrides: []byte(DefaultWhereElement_PrefixOverrides),
				suffixOverrides: []byte(v.Propertys["suffixOverrides"]),
				childs:          []Node{},
			}
			if v.ElementItems != nil && len(v.ElementItems) > 0 {
				var childNodes = it.ParserNodes(v.ElementItems)
				n.childs = append(n.childs, childNodes...)
			} else {
				n.childs = nil
			}
			node = &n
			break
		case "bind":
			n := NodeBind{
				t:      NBind,
				value:  v.Propertys["value"],
				name:   v.Propertys["name"],
				holder: &it.holder,
			}
			node = &n
		}
		nodes = append(nodes, node)
	}
	return nodes
}

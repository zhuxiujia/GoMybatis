package GoMybatis

type NodeParser struct {
}

func (it NodeParser) ParserNodes(mapperXml []ElementItem) []SqlNode {
	var nodes = []SqlNode{}
	for _, v := range mapperXml {
		var node SqlNode

		switch v.ElementType {
		case "string":
			n := StringNode{
				value: v.DataString,
				t:     NString,
			}
			node = &n
			break
		case "if":
			n := IfNode{
				t:      NIf,
				test:   v.Propertys["test"],
				childs: []SqlNode{},
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
			n := TrimNode{
				t:               NTrim,
				prefix:          []byte(v.Propertys["prefix"]),
				suffix:          []byte(v.Propertys["suffix"]),
				prefixOverrides: []byte(v.Propertys["prefixOverrides"]),
				suffixOverrides: []byte(v.Propertys["suffixOverrides"]),
				childs:          []SqlNode{},
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
			n := TrimNode{
				t:      NTrim,
				childs: []SqlNode{},

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
			n := ForEachNode{
				t:          NForEach,
				childs:     []SqlNode{},
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
			n := ChooseNode{
				t:         NChoose,
				whenNodes: []SqlNode{},
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
		case "otherwise":
			n := OtherwiseNode{
				t:      NOtherwise,
				childs: []SqlNode{},
			}
			if v.ElementItems != nil && len(v.ElementItems) > 0 {
				var childNodes = it.ParserNodes(v.ElementItems)
				n.childs = append(n.childs, childNodes...)
			} else {
				n.childs = nil
			}
			node = &n
			break
		case "when":
			n := WhenNode{
				t:      NOtherwise,
				childs: []SqlNode{},
				test:   v.Propertys["test"],
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
			n := TrimNode{
				t:               NTrim,
				prefix:          []byte(DefaultWhereElement_Prefix),
				suffix:          []byte(v.Propertys["suffix"]),
				prefixOverrides: []byte(DefaultWhereElement_PrefixOverrides),
				suffixOverrides: []byte(v.Propertys["suffixOverrides"]),
				childs:          []SqlNode{},
			}
			if v.ElementItems != nil && len(v.ElementItems) > 0 {
				var childNodes = it.ParserNodes(v.ElementItems)
				n.childs = append(n.childs, childNodes...)
			} else {
				n.childs = nil
			}
			node = &n
			break
		}
		nodes = append(nodes, node)
	}
	return nodes
}

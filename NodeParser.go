package GoMybatis

type NodeParser struct {
}

func (it NodeParser) ParserNodes(mapperXml []ElementItem) []SqlNode {
	var nodes = []SqlNode{}
	for _, v := range mapperXml {
		var node SqlNode
		if v.ElementType == "string" {
			n := StringNode{
				value: v.DataString,
				t:     NString,
			}
			node = n
		} else if v.ElementType == "if" {
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
			node = n
		}
		nodes = append(nodes, node)
	}
	return nodes
}

package GoExpress

import (
	"fmt"
	"strconv"
	"testing"
)

func TestNode_Run(t *testing.T) {
	var express = "1"

	//express = "1 + 2 > 3 + 6"

	//express = "1 + 2 != nil"

	var opts = ParserOperators(express)

	fmt.Println(opts)

	var nodes []node
	for _, v := range opts {
		var node node
		if node == nil && v == "nil" {
			var inode = NilNode{}
			node = inode
		}

		if isOperatorsAction(v) {
			var optNode = OptNode{
				value: v,
			}
			node = optNode
		}

		var i, e = strconv.ParseInt(v, 0, 64)
		if node == nil && e == nil {
			var inode = IntNode{
				value: int64(i),
			}
			node = inode
		}
		u, _ := strconv.ParseUint(v, 0, 64)
		if node == nil && e == nil {
			var inode = UIntNode{
				value: u,
			}
			node = inode
		}
		f, e := strconv.ParseFloat(v, 64)
		if node == nil && e == nil {
			var inode = FloatNode{
				value: f,
			}
			node = inode
		}
		b, e := strconv.ParseBool(v)
		if node == nil && e == nil {
			var inode = BoolNode{
				value: b,
			}
			node = inode
		}
		if node == nil {
			var argNode = ArgNode{
				value: v,
			}
			node = argNode
		}
		if node == nil {
			panic("uncheck opt " + v)
		}
		nodes = append(nodes, node)
	}

	fmt.Println(nodes)

	for _, v := range priorityArray {
		findReplaceOpt(v, &nodes)
	}

	fmt.Println(nodes)
}

func findReplaceOpt(operator Operator, nodearg *[]node) {
	var nodes = *nodearg
	for nIndex, n := range nodes {
		if n.Type() == NOpt {
			if nIndex == 0 || (nIndex+1) == len(nodes) {
				panic("expr not true!")
			}
			if nIndex-1 > 0 && nodes[nIndex-1].Type() == NOpt {
				panic("expr not true!")
			}
			if nIndex < len(nodes) && nodes[nIndex+1].Type() == NOpt {
				panic("expr not true!")
			}
			var opt = n.(OptNode)
			if opt.value != operator {
				continue
			}

			var newNode node
			if isOperatorsAction(opt.value) {
				newNode = CalculationNode{
					left:  nodes[nIndex-1],
					right: nodes[nIndex+1],
					opt:   opt.value,
				}
			} else {
				newNode = EqualNode{
					left:  nodes[nIndex-1],
					right: nodes[nIndex+1],
					opt:   opt.value,
				}
			}
			var newNodes []node
			newNodes = append(nodes[:nIndex-1], newNode)
			newNodes = append(newNodes, nodes[nIndex+2:]...)

			if haveOpt(newNodes) {
				findReplaceOpt(operator, &newNodes)
			}
			*nodearg = newNodes
			break
		}
	}
}

func haveOpt(nodes []node) bool {
	for _, v := range nodes {
		if v.Type() == NOpt {
			return true
		}
	}
	return false
}

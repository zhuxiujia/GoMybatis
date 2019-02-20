package GoExpress

import (
	"errors"
	"strconv"
	"strings"
)

//操作符
type Operator = string

const (
	//计算操作符
	Add    Operator = "+"
	Reduce Operator = "-"
	Ride   Operator = "*"
	Divide Operator = "/"

	//比较操作符
	And       Operator = "&&"
	Or        Operator = "||"
	Equal     Operator = "=="
	UnEqual   Operator = "!="
	Less      Operator = "<"
	LessEqual Operator = "<="
	More      Operator = ">"
	MoreEqual Operator = ">="
)

//乘除优先于加减 计算优于比较,
var priorityArray = []Operator{Ride, Divide, Add, Reduce,
	LessEqual, Less, MoreEqual, More,
	UnEqual, Equal, And, Or}

//操作符优先级
var priorityMap = map[Operator]int{}

func init() {
	for k, v := range priorityArray {
		priorityMap[v] = k
	}
}

func Parser(express string) (Node, error) {
	var opts = ParserOperators(express)

	var nodes []Node
	for _, v := range opts {
		var node = parserNode(v)
		nodes = append(nodes, node)
	}

	for _, v := range priorityArray {
		var e = findReplaceOpt(v, &nodes)
		if e != nil {
			return nil, e
		}
	}
	if len(nodes) == 0 || nodes[0] == nil {
		return nil, errors.New("parser node fail!")
	}
	return nodes[0], nil
}

func parserNode(v Operator) Node {
	var node Node
	if v == "nil" {
		var inode = NilNode{}
		node = inode
	}

	if isOperatorsAction(v) {
		var optNode = OptNode{
			value: v,
		}
		node = optNode
	}

	i, e := strconv.ParseInt(v, 0, 64)
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
	e = nil
	if node == nil && e == nil &&
		strings.Index(v, "'") == 0 && strings.LastIndex(v, "'") == (len(v)-1) {
		var inode = StringNode{
			value: string([]byte(v)[1 : len(v)-1]),
		}
		node = inode
	}
	e = nil

	if node == nil {
		var argNode = ArgNode{
			value: v,
		}
		node = argNode
	}
	if node == nil {
		panic("uncheck opt " + v)
	}
	return node
}
func findReplaceOpt(operator Operator, nodearg *[]Node) error {
	var nodes = *nodearg
	for nIndex, n := range nodes {
		if n.Type() == NOpt {
			if nIndex == 0 || (nIndex+1) == len(nodes) {
				return errors.New("expr operator" + operator + " left or right not have value!")
			}
			if nIndex-1 > 0 && nodes[nIndex-1].Type() == NOpt {
				return errors.New("expr same operator can not have more than 2!")
			}
			if nIndex < len(nodes) && nodes[nIndex+1].Type() == NOpt {
				return errors.New("expr not true!")
			}
			var opt = n.(OptNode)
			if opt.value != operator {
				continue
			}

			var newNode = BinaryNode{
				left:  nodes[nIndex-1],
				right: nodes[nIndex+1],
				opt:   opt.value,
			}
			var newNodes []Node
			newNodes = append(nodes[:nIndex-1], newNode)
			newNodes = append(newNodes, nodes[nIndex+2:]...)

			if haveOpt(newNodes) {
				findReplaceOpt(operator, &newNodes)
			}
			*nodearg = newNodes
			break
		}
	}

	return nil
}

func haveOpt(nodes []Node) bool {
	for _, v := range nodes {
		if v.Type() == NOpt {
			return true
		}
	}
	return false
}

func ParserOperators(express string) []Operator {
	express = strings.Replace(express, "nil", " nil ", -1)
	express = strings.Replace(express, Add, " "+Add+" ", -1)
	express = strings.Replace(express, Reduce, " "+Reduce+" ", -1)
	express = strings.Replace(express, Ride, " "+Ride+" ", -1)
	express = strings.Replace(express, Divide, " "+Divide+" ", -1)
	express = strings.Replace(express, And, " "+And+" ", -1)
	express = strings.Replace(express, Or, " "+Or+" ", -1)
	express = strings.Replace(express, UnEqual, " "+UnEqual+" ", -1)
	express = strings.Replace(express, Equal, " "+Equal+" ", -1)
	express = strings.Replace(express, LessEqual, " "+LessEqual+" ", -1)
	express = strings.Replace(express, Less, " "+Less+" ", -1)
	express = strings.Replace(express, MoreEqual, " "+MoreEqual+" ", -1)
	express = strings.Replace(express, More, " "+More+" ", -1)

	var newResult []string
	var results = strings.Split(express, " ")
	for _, v := range results {
		if v != " " && v != "" {
			newResult = append(newResult, v)
		}
	}
	return newResult
}

func isOperatorsAction(arg string) bool {
	//计算操作符
	if arg == Add ||
		arg == Reduce ||
		arg == Ride ||
		arg == Divide ||
		//比较操作符
		arg == And ||
		arg == Or ||
		arg == Equal ||
		arg == UnEqual ||
		arg == Less ||
		arg == LessEqual ||
		arg == More ||
		arg == MoreEqual {
		return true
	}
	return false

}

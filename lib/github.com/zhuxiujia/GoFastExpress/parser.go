package GoFastExpress

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

	Nil  Operator = "nil"
	Null Operator = "null"
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
		var node, err = parserNode(express, v)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}
	//check epress
	var err = checkeNodes(express, nodes)
	if err != nil {
		return nil, err
	}
	for _, v := range priorityArray {
		var e = findReplaceOpt(express, v, &nodes)
		if e != nil {
			return nil, e
		}
	}
	if len(nodes) == 0 || nodes[0] == nil {
		return nil, errors.New("parser node fail!")
	}
	return nodes[0], nil
}

func checkeNodes(express string, nodes []Node) error {
	var nodesLen = len(nodes)
	for nIndex, n := range nodes {
		if n.Type() == NOpt {
			var after Node
			var befor Node

			if nIndex > 0 {
				befor = nodes[nIndex-1]
			}
			if nIndex < (nodesLen - 1) {
				after = nodes[nIndex+1]
			}
			if after != nil && after.Type() == NOpt {
				return errors.New("express have more than 2 opt!express=" + express)
			}
			if befor != nil && befor.Type() == NOpt {
				return errors.New("express have more than 2 opt!express=" + express)
			}
		}
	}
	return nil
}

func parserNode(express string, v Operator) (Node, error) {
	if v == Nil || v == Null {
		var inode = NilNode{
			t: NNil,
		}
		return inode, nil
	}
	if v == "=" {
		return nil, errors.New("find not support opt = '=',express=" + express)
	}
	if isOperatorsAction(v) {
		var optNode = OptNode{
			value: v,
			t:     NOpt,
		}
		return optNode, nil
	}

	i, e := strconv.ParseInt(v, 0, 64)
	if e == nil {
		var inode = IntNode{
			value: int64(i),
			t:     NInt,
		}
		return inode, nil
	}
	u, _ := strconv.ParseUint(v, 0, 64)
	if e == nil {
		var inode = UIntNode{
			value: u,
			t:     NUInt,
		}
		return inode, nil
	}
	f, e := strconv.ParseFloat(v, 64)
	if e == nil {
		var inode = FloatNode{
			value: f,
			t:     NFloat,
		}
		return inode, nil
	}
	b, e := strconv.ParseBool(v)
	if e == nil {
		var inode = BoolNode{
			value: b,
			t:     NBool,
		}
		return inode, nil
	}
	if strings.Index(v, "'") == 0 && strings.LastIndex(v, "'") == (len(v)-1) {
		var inode = StringNode{
			value: string([]byte(v)[1 : len(v)-1]),
			t:     NString,
		}
		return inode, nil
	}
	e = nil
	if isOperatorsAction(v) {
		var optNode = OptNode{
			value: v,
			t:     NOpt,
		}
		return optNode, nil
	}
	var argNode = ArgNode{
		value: v,
		t:     NArg,
	}
	return argNode, nil
}
func findReplaceOpt(express string, operator Operator, nodearg *[]Node) error {
	var nodes = *nodearg
	for nIndex, n := range nodes {
		if n.Type() == NOpt {
			var opt = n.(OptNode)
			if opt.value != operator {
				continue
			}
			var newNode = BinaryNode{
				left:  nodes[nIndex-1],
				right: nodes[nIndex+1],
				opt:   opt.value,
				t:     NBinary,
			}
			var newNodes []Node
			newNodes = append(nodes[:nIndex-1], newNode)
			newNodes = append(newNodes, nodes[nIndex+2:]...)

			if haveOpt(newNodes) {
				findReplaceOpt(express, operator, &newNodes)
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
	express = strings.Replace(express, Null, " "+Null+" ", -1)
	express = strings.Replace(express, Nil, " "+Nil+" ", -1)
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

	express = strings.Replace(express, "! =", " "+UnEqual+" ", -1)
	express = strings.Replace(express, "= =", " "+Equal+" ", -1)
	express = strings.Replace(express, "< =", " "+LessEqual+" ", -1)
	express = strings.Replace(express, "> =", " "+MoreEqual+" ", -1)
	express = strings.Replace(express, "& &", " "+And+" ", -1)
	express = strings.Replace(express, "| |", " "+Or+" ", -1)

	var newResult []string
	var results = strings.Split(express, " ")
	for _, v := range results {
		if v != "" && v != " " {
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

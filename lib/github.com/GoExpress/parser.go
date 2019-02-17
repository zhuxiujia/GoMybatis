package GoExpress

import (
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
	Equal     Operator = "="
	UnEqual   Operator = "!="
	Less      Operator = "<"
	LessEqual Operator = "<="
	More      Operator = ">"
	MoreEqual Operator = ">="
)

//节点类型
type NodeType = int

const (
	Calculation NodeType = iota //计算节点
	EqualValue                  //比较节点
	Value                       //值节点
)

//节点
type Node struct {
	NodeType NodeType
	Operator Operator

	LeftOpt  Operator
	RightOpt Operator

	Left      interface{}
	Right     interface{}
	NodeValue interface{} //节点执行结果
}

func (it *Node) Run() (interface{}, error) {
	if it.NodeType == Calculation {
		return Eval(it.Operator, it.Left, it.Right)
	} else {
		return Eval(it.Operator, it.Left, it.Right)
	}
}

//解析
func Parser(express string) []Node {
	var nodes []Node
	var newResult = ParserOperators(express)
	var tempNode = Node{}
	for _, v := range newResult {
		var isOpt = isOperatorsAction(v)
		if isOpt {
			tempNode.Operator = v
		} else {
			var result interface{}
			if strings.Index(v, "'") == 0 && strings.LastIndex(v, "'") == (len(v)-1) {
				var bytes = []byte(v)[1 : len(v)-1]
				result = string(bytes)
			} else {
				var i, e = strconv.ParseInt(v, 0, 64)
				if e == nil {
					result = i
				}
				u, _ := strconv.ParseUint(v, 0, 64)
				if e == nil {
					result = u
				}
				f, e := strconv.ParseFloat(v, 64)
				if e == nil {
					result = f
				}
				b, e := strconv.ParseBool(v)
				if e == nil {
					result = b
				}
			}
			tempNode.NodeValue = result
			if tempNode.LeftOpt != "" {
				tempNode.RightOpt = v
			}
			tempNode.LeftOpt = v
		}
		if tempNode.LeftOpt != "" && tempNode.RightOpt != "" {
			nodes = append(nodes, tempNode)
			tempNode = Node{}
		}
	}
	return nodes
}

func ParserOperators(express string) []Operator {
	express = strings.Replace(express, Equal, " "+Equal+" ", -1)
	express = strings.Replace(express, Reduce, " "+Reduce+" ", -1)
	express = strings.Replace(express, Ride, " "+Ride+" ", -1)
	express = strings.Replace(express, Divide, " "+Divide+" ", -1)
	express = strings.Replace(express, And, " "+And+" ", -1)
	express = strings.Replace(express, Or, " "+Or+" ", -1)
	express = strings.Replace(express, Equal, " "+Equal+" ", -1)
	express = strings.Replace(express, UnEqual, " "+UnEqual+" ", -1)
	express = strings.Replace(express, Less, " "+Less+" ", -1)
	express = strings.Replace(express, LessEqual, " "+LessEqual+" ", -1)
	express = strings.Replace(express, More, " "+More+" ", -1)
	express = strings.Replace(express, MoreEqual, " "+MoreEqual+" ", -1)
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

package GoFastExpress

import (
	"errors"
	"go/scanner"
	"go/token"
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

var NotSupportOptMap = map[string]bool{
	"=": true,
	"!": true,
	"@": true,
	"#": true,
	"$": true,
	"^": true,
	"&": true,
	"(": true,
	")": true,
	"`": true,
}

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
	if NotSupportOptMap[v] {
		return nil, errors.New("find not support opt = '" + v + "',express=" + express)
	}
	if isOperatorsAction(v) {
		var optNode = OptNode{
			value: v,
			t:     NOpt,
		}
		return optNode, nil
	}
	if v == "true" || v == "false" {
		b, e := strconv.ParseBool(v)
		if e == nil {
			var inode = BoolNode{
				value: b,
				t:     NBool,
			}
			return inode, nil
		}
	}
	if strings.Index(v, "'") == 0 && strings.LastIndex(v, "'") == (len(v)-1) {
		var inode = StringNode{
			value: string([]byte(v)[1 : len(v)-1]),
			t:     NString,
		}
		return inode, nil
	}
	if strings.Index(v, "`") == 0 && strings.LastIndex(v, "`") == (len(v)-1) {
		var inode = StringNode{
			value: string([]byte(v)[1 : len(v)-1]),
			t:     NString,
		}
		return inode, nil
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
	e = nil

	var values=strings.Split(v,".")
	var argNode = ArgNode{
		value: v,
		values:values,
		valuesLen: len(values),
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
	var newResult []string
	src := []byte(express)
	var s scanner.Scanner
	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(src))
	s.Init(file, src, nil, 0)
	var lastToken token.Token
	for {
		_, tok, lit := s.Scan()
		if tok == token.EOF || lit == "\n" {
			break
		}
		//fmt.Printf("%-6s%-8s%q\n", fset.Position(pos), tok, lit)
		var s = toStr(lit, tok)
		if lit == "" && tok != token.ILLEGAL {
			lastToken = tok
		}
		if tok == token.PERIOD || lastToken == token.PERIOD {
			//append to last token
			newResult[len(newResult)-1] = newResult[len(newResult)-1] + s
			continue
		}
		newResult = append(newResult, s)
	}
	return newResult
}

func toStr(lit string, tok token.Token) string {
	if lit == "" {
		return tok.String()
	} else {
		return lit
	}
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

package GoMybatis

import (
	"bytes"
)

//sql构建抽象语法树节点
type Node interface {
	Type() NodeType
	Eval(env map[string]interface{}) ([]byte, error)
}

//执行子所有节点
func DoChildNodes(childNodes []Node, env map[string]interface{}) ([]byte, error) {
	if childNodes == nil {
		return nil, nil
	}
	var sql bytes.Buffer
	for _, v := range childNodes {
		var r, e = v.Eval(env)
		if e != nil {
			return nil, e
		}
		if r != nil {
			sql.Write(r)
		}
	}
	return sql.Bytes(), nil
}

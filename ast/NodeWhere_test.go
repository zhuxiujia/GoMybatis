package ast

import (
	"strings"
	"testing"
)

func TestNodeWhere_Eval(t *testing.T) {
	var w = NodeWhere{
		t: NWhere,
	}
	w.childs = []Node{&NodeString{
		value:               "and select *",
		t:                   NString,
		expressMap:          []string{},
		noConvertExpressMap: []string{},
		holder:              &NodeConfigHolder{},
	}}

	var s, e = w.Eval(nil, nil, nil)
	if e != nil {
		t.Fatal(e)
	}
	if !strings.HasPrefix(string(s), " WHERE") {
		t.Fatal("where node fail!")
	}
	println(string(s))
}

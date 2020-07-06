package plugin

import (
	"encoding/json"
	"testing"
)

type A struct {
	Name string `json:"name"`
}

func Test_page(t *testing.T) {
	var p = Page{[]json.RawMessage{}}
	var a = A{
		Name: "xiao ming",
	}
	var ajs, _ = json.Marshal(&a)
	p.Content = append(p.Content, ajs)
	var s, _ = json.Marshal(&p)
	var js = string(s)
	println("mashal:", js)

	var np = Page{}
	json.Unmarshal(s, &np)

	var contents = []A{}
	np.ParserContent(&contents)
	println("contents", contents)
}

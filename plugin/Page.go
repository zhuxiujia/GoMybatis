package plugin

import "encoding/json"

//TODO 分页插件

type Page struct {
	Content []json.RawMessage `json:"content"`
}

// parser content to type
func (it *Page) ParserContent(result interface{}) error {
	var jsBytes, _ = json.Marshal(it.Content)
	return json.Unmarshal(jsBytes, &result)
}

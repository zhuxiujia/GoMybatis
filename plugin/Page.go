package plugin

import "encoding/json"

//TODO 分页插件,在go 2.0支持泛型后再计划改为泛型Page
type Page struct {
	Content []json.RawMessage `json:"content"`
}

// parser content to type
func (it *Page) GetContent(contentArray interface{}) error {
	var jsBytes, e = json.Marshal(it.Content)
	if e != nil {
		return e
	}
	return json.Unmarshal(jsBytes, &contentArray)
}

// parser content to type
func (it *Page) GetContentData(index int, content interface{}) error {
	return json.Unmarshal(it.Content[index], &content)
}

//append one content into content array
func (it *Page) AppendContentData(content interface{}) error {
	var jsBytes, e = json.Marshal(content)
	if e != nil {
		return e
	}
	it.Content = append(it.Content, jsBytes)
	return nil
}

// contentArray interface must be array
func (it *Page) SetContent(contentArray interface{}) error {
	var jsBytes, e = json.Marshal(contentArray)
	if e != nil {
		return e
	}
	var data = []json.RawMessage{}
	e = json.Unmarshal(jsBytes, &data)
	if e != nil {
		return e
	}
	it.Content = data
	return nil
}

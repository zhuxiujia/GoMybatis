package GoMybatis

type TempleteDecoder interface {
	Decode(mapper *MapperXml) error
	DecodeTree(tree map[string]*MapperXml) error
}

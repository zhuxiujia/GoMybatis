package GoMybatis

type TempleteDecoder interface {
	Decode(mapper *MapperXml) error
}

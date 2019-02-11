package GoMybatis

import "reflect"

type TempleteDecoder interface {
	DecodeTree(tree map[string]*MapperXml, beanType reflect.Type) error
}

package GoMybatis

import (
	"GoMybatis/lib/github.com/beevik/etree"
	"reflect"
)

type TempleteDecoder interface {
	DecodeTree(tree map[string]etree.Token, beanType reflect.Type) error
}

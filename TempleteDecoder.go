package GoMybatis

import (
	"github.com/zhuxiujia/GoMybatis/lib/github.com/beevik/etree"
	"reflect"
)

type TempleteDecoder interface {
	DecodeTree(tree map[string]*etree.Element, beanType reflect.Type) error
}

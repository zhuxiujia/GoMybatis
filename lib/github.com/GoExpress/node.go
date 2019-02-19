package GoExpress

type nodeType int

const (
	NArg    nodeType = iota
	NString          //string 节点
	NFloat           //float节点
	NInt             //int 节点
	NUInt            //uint节点
	NBool            //bool节点
	NNil             //空节点
	NBinary          //二元计算节点
	NOpt             //操作符节点
)

//抽象语法树节点
type Node interface {
	Type() nodeType
	Eval(env interface{}) (interface{}, error)
}

type OptNode struct {
	value Operator
}

func (OptNode) Type() nodeType {
	return NOpt
}
func (it OptNode) IsCalculationOperator() bool {
	//计算操作符
	if it.value == Add ||
		it.value == Reduce ||
		it.value == Ride ||
		it.value == Divide {
		return true
	}
	return false

}

func (it OptNode) Eval(env interface{}) (interface{}, error) {
	return it.value, nil
}

//参数节点
type ArgNode struct {
	value string
}

func (ArgNode) Type() nodeType {
	return NArg
}

func (it ArgNode) Eval(env interface{}) (interface{}, error) {
	//TODO do arg
	return EvalTake(it.value, env)
}

//值节点
type StringNode struct {
	value string
}

func (StringNode) Type() nodeType {
	return NString
}

func (it StringNode) Eval(env interface{}) (interface{}, error) {
	return it.value, nil
}

//值节点
type FloatNode struct {
	value float64
}

func (FloatNode) Type() nodeType {
	return NFloat
}

func (it FloatNode) Eval(env interface{}) (interface{}, error) {
	return it.value, nil
}

//值节点
type IntNode struct {
	value int64
}

func (IntNode) Type() nodeType {
	return NInt
}

func (it IntNode) Eval(env interface{}) (interface{}, error) {
	return it.value, nil
}

type UIntNode struct {
	value uint64
}

func (UIntNode) Type() nodeType {
	return NUInt
}

func (it UIntNode) Eval(env interface{}) (interface{}, error) {
	return it.value, nil
}

//值节点
type BoolNode struct {
	value bool
}

func (BoolNode) Type() nodeType {
	return NBool
}

func (it BoolNode) Eval(env interface{}) (interface{}, error) {
	return it.value, nil
}

//空节点
type NilNode struct {
}

func (NilNode) Type() nodeType {
	return NNil
}

func (NilNode) Eval(env interface{}) (interface{}, error) {
	return nil, nil
}

//计算节点
type BinaryNode struct {
	left  Node
	right Node
	opt   Operator
}

func (BinaryNode) Type() nodeType {
	return NBinary
}

func (it BinaryNode) Eval(env interface{}) (interface{}, error) {
	var left interface{}
	var right interface{}
	var e error
	if it.left != nil {
		left, e = it.left.Eval(env)
		if e != nil {
			return nil, e
		}
	}
	if it.right != nil {
		right, e = it.right.Eval(env)
		if e != nil {
			return nil, e
		}
	}
	return Eval(it.opt, left, right)
}

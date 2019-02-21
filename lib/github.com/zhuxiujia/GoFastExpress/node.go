package GoFastExpress

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

func (it nodeType) ToString() string {
	switch it {
	case NArg:
		return "NArg"
	case NString:
		return "NString"
	case NFloat:
		return "NFloat"
	case NInt:
		return "NInt"
	case NUInt:
		return "NUInt"
	case NBool:
		return "NBool"
	case NNil:
		return "NNil"
	case NBinary:
		return "NBinary"
	case NOpt:
		return "NOpt"
	}
	return "Unknow"
}

//抽象语法树节点
type Node interface {
	Type() nodeType
	Eval(env interface{}) (interface{}, error)
}

type OptNode struct {
	value Operator
	t     nodeType
}

func (it OptNode) Type() nodeType {
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
	t     nodeType
}

func (it ArgNode) Type() nodeType {
	return NArg
}

func (it ArgNode) Eval(env interface{}) (interface{}, error) {
	//TODO do arg
	return EvalTake(it.value, env)
}

//值节点
type StringNode struct {
	value string
	t     nodeType
}

func (it StringNode) Type() nodeType {
	return NString
}

func (it StringNode) Eval(env interface{}) (interface{}, error) {
	return it.value, nil
}

//值节点
type FloatNode struct {
	value float64
	t     nodeType
}

func (it FloatNode) Type() nodeType {
	return NFloat
}

func (it FloatNode) Eval(env interface{}) (interface{}, error) {
	return it.value, nil
}

//值节点
type IntNode struct {
	value int64
	t     nodeType
}

func (it IntNode) Type() nodeType {
	return NInt
}

func (it IntNode) Eval(env interface{}) (interface{}, error) {
	return it.value, nil
}

type UIntNode struct {
	value uint64
	t     nodeType
}

func (it UIntNode) Type() nodeType {
	return NUInt
}

func (it UIntNode) Eval(env interface{}) (interface{}, error) {
	return it.value, nil
}

//值节点
type BoolNode struct {
	value bool
	t     nodeType
}

func (it BoolNode) Type() nodeType {
	return NBool
}

func (it BoolNode) Eval(env interface{}) (interface{}, error) {
	return it.value, nil
}

//空节点
type NilNode struct {
	t nodeType
}

func (it NilNode) Type() nodeType {
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
	t     nodeType
}

func (it BinaryNode) Type() nodeType {
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

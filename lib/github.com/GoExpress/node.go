package GoExpress

type nodeType int

const (
	NArg         nodeType = iota
	NString               //string 节点
	NFloat                //float节点
	NInt                  //int 节点
	NUInt                 //uint节点
	NBool                 //bool节点
	NNil                  //空节点
	NEqual                //比较节点
	NCalculation          //计算节点
	NOpt                  //操作符节点
)

//抽象语法树节点
type node interface {
	Type() nodeType
	Eval(operator Operator, nexNode Node) (interface{}, error)
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

func (OptNode) Eval(operator Operator, nexNode Node) (interface{}, error) {
	return nil, nil
}

//参数节点
type ArgNode struct {
	value string
}

func (ArgNode) Type() nodeType {
	return NArg
}

func (ArgNode) Eval(operator Operator, nexNode Node) (interface{}, error) {
	return nil, nil
}

//值节点
type StringNode struct {
	value string
}

func (StringNode) Type() nodeType {
	return NString
}

func (StringNode) Eval(operator Operator, nexNode Node) (interface{}, error) {
	return nil, nil
}

//值节点
type FloatNode struct {
	value float64
}

func (FloatNode) Type() nodeType {
	return NFloat
}

func (FloatNode) Eval(operator Operator, nexNode Node) (interface{}, error) {
	return nil, nil
}

//值节点
type IntNode struct {
	value int64
}

func (IntNode) Type() nodeType {
	return NInt
}

func (IntNode) Eval(operator Operator, nexNode Node) (interface{}, error) {
	return nil, nil
}

type UIntNode struct {
	value uint64
}

func (UIntNode) Type() nodeType {
	return NUInt
}

func (UIntNode) Eval(operator Operator, nexNode Node) (interface{}, error) {
	return nil, nil
}

//值节点
type BoolNode struct {
	value bool
}

func (BoolNode) Type() nodeType {
	return NBool
}

func (BoolNode) Eval(operator Operator, nexNode Node) (interface{}, error) {
	return nil, nil
}

//空节点
type NilNode struct {
}

func (NilNode) Type() nodeType {
	return NNil
}

func (NilNode) Eval(operator Operator, nexNode Node) (interface{}, error) {
	return nil, nil
}

/**
组合节点
*/
//比较节点
type EqualNode struct {
	left  node
	right node
	opt   Operator
}

func (EqualNode) Type() nodeType {
	return NEqual
}

func (EqualNode) Eval(operator Operator, nexNode Node) (interface{}, error) {
	return nil, nil
}

//计算节点
type CalculationNode struct {
	left  node
	right node
	opt   Operator
}

func (CalculationNode) Type() nodeType {
	return NCalculation
}

func (CalculationNode) Eval(operator Operator, nexNode Node) (interface{}, error) {
	return nil, nil
}

// IntNode opt IntNode Opt IntNode Opt IntNode
// => CalculationNode

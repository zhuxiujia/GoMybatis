package GoExpress

//抽象语法数节点
type node interface {
	Eval(env interface{}) (interface{}, error)
}

//参数节点
type ArgNode struct {
	value string
}

//值节点
type StringNode struct {
	value string
}

//值节点
type FloatNode struct {
	value float64
}

//值节点
type IntNode struct {
	value int
}

//空节点
type NilNode struct {
}

//比较节点
type EqualNode struct {
	left  node
	right node
}

//计算节点
type CalculationNode struct {
	left  node
	right node
}

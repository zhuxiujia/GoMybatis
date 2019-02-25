package GoMybatis

type SqlNodeType int

const (
	NArg    SqlNodeType = iota
	NString             //string 节点
	NNil                //空节点
	NBinary             //二元计算节点
	NOpt                //操作符节点

	NIf
	NTrim
	NSet
	NForEach
	NChoose
	NOtherwise
)

func (it SqlNodeType) ToString() string {
	switch it {
	case NString:
		return "NString"
	case NNil:
		return "NNil"
	case NBinary:
		return "NBinary"
	case NOpt:
		return "NOpt"
	}
	return "Unknow"
}

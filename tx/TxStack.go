package tx

import "reflect"

//事务栈
type TxStack struct {
	i    int
	data []reflect.StructField //方法队列
}

func (it TxStack) New() TxStack {
	return TxStack{
		data: []reflect.StructField{},
		i:    0,
	}
}

func (s *TxStack) Push(k reflect.StructField) {
	s.data = append(s.data, k)
	s.i++
}

func (s *TxStack) Pop() (ret reflect.StructField) {
	s.i--
	ret = s.data[s.i]
	s.data = s.data[0:s.i]
	return
}

func (s *TxStack) Len() int {
	return s.i
}

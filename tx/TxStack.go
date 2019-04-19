package tx

//事务栈
type TxStack struct {
	i    int
	data []string
}

func (it TxStack) New() TxStack {
	return TxStack{
		data: []string{},
		i:    0,
	}
}

func (s *TxStack) Push(k string) {
	s.data = append(s.data, k)
	s.i++
}

func (s *TxStack) Pop() (ret string) {
	s.i--
	ret = s.data[s.i]
	s.data = s.data[0:s.i]
	return
}

package tx

//事务栈
type TxStack struct {
	i    int
	data []string
}

func (s *TxStack) Push(k string) {
	s.data[s.i] = k
	s.i++
}

func (s *TxStack) Pop() (ret string) {
	s.i--
	ret = s.data[s.i]
	return
}

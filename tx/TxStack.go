package tx

import "database/sql"

type TxStack struct {
	i    int
	data []*sql.Tx //队列
}

func (it TxStack) New() TxStack {
	return TxStack{
		data: []*sql.Tx{},
		i:    0,
	}
}

func (s *TxStack) Push(k *sql.Tx) {
	s.data = append(s.data, k)
	s.i++
}

func (s *TxStack) Pop() *sql.Tx {
	if s.i == 0 {
		return nil
	}
	s.i--
	var ret = s.data[s.i]
	s.data = s.data[0:s.i]
	return ret
}
func (s *TxStack) First() *sql.Tx {
	if s.i == 0 {
		return nil
	}
	var ret = s.data[0]
	return ret
}
func (s *TxStack) Last() *sql.Tx {
	if s.i == 0 {
		return nil
	}
	var ret = s.data[s.i-1]
	return ret
}

func (s *TxStack) Len() int {
	return s.i
}

func (s *TxStack) HaveTx() bool {
	return s.Len() > 0
}

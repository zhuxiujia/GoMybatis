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

func (s *TxStack) Pop() (ret *sql.Tx) {
	s.i--
	ret = s.data[s.i]
	s.data = s.data[0:s.i]
	return
}

func (s *TxStack) Len() int {
	return s.i
}


package tx

import "database/sql"

type TxStack struct {
	i            int
	data         []*sql.Tx      //队列
	propagations []*Propagation //队列
}

func (it TxStack) New() TxStack {
	return TxStack{
		data:         []*sql.Tx{},
		propagations: []*Propagation{},
		i:            0,
	}
}

func (s *TxStack) Push(k *sql.Tx, p *Propagation) {
	s.data = append(s.data, k)
	s.propagations = append(s.propagations, p)
	s.i++
}

func (s *TxStack) Pop() (*sql.Tx, *Propagation) {
	if s.i == 0 {
		return nil, nil
	}
	s.i--
	var ret = s.data[s.i]
	s.data = s.data[0:s.i]

	var p = s.propagations[s.i]
	s.propagations = s.propagations[0:s.i]
	return ret, p
}
func (s *TxStack) First() (*sql.Tx, *Propagation) {
	if s.i == 0 {
		return nil, nil
	}
	var ret = s.data[0]
	var p = s.propagations[0]
	return ret, p
}
func (s *TxStack) Last() (*sql.Tx, *Propagation) {
	if s.i == 0 {
		return nil, nil
	}
	var ret = s.data[s.i-1]
	var p = s.propagations[s.i-1]
	return ret, p
}

func (s *TxStack) Len() int {
	return s.i
}

func (s *TxStack) HaveTx() bool {
	return s.Len() > 0
}

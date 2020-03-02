package tx

import "reflect"

//session map是协程安全的
//此处无需处理并发，因为都是单协程访问
type StructField struct {
	i    int
	data []reflect.StructField //方法队列
}

func (it StructField) New() StructField {
	return StructField{
		data: []reflect.StructField{},
		i:    0,
	}
}

func (s *StructField) Push(k reflect.StructField) {
	s.data = append(s.data, k)
	s.i++
}

func (s *StructField) Pop() (ret reflect.StructField) {
	s.i--
	ret = s.data[s.i]
	s.data = s.data[0:s.i]
	return
}

func (s *StructField) Len() int {
	return s.i
}

package GoMybatis

import "log"

type LogStandard struct {
	PrintlnFunc func(v []byte)
}

func (this *LogStandard) Println(v []byte) {
	if this.PrintlnFunc != nil {
		this.PrintlnFunc(v)
	} else {
		log.Println(string(v))
	}
}

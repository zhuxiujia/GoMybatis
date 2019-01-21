package GoMybatis

import "log"

type LogStandard struct {
	PrintlnFunc func(v ...interface{})
}

func (this *LogStandard) Println(v ...interface{}) {
	if this.PrintlnFunc != nil {
		this.PrintlnFunc(v)
	} else {
		var infos = v
		log.Println(infos)
	}
}

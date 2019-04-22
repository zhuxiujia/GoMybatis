package utils

import (
	"github.com/zhuxiujia/GoMybatis/lib/github.com/goroutineid"
)

func GoroutineID() int64 {
	return goroutineid.GetGoID()
}

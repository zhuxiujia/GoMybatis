package utils

import (
	"fmt"
	"reflect"
	"time"
)

var c_FLOAT64_DEFAULT float64
var Float64Type = reflect.TypeOf(c_FLOAT64_DEFAULT)

//调用次方法 必须加上 defer 关键字！
func CountMethodUseTime(now time.Time, info string, duration time.Duration) {
	var end = time.Now()
	var durationName = DurationToString(duration)
	fmt.Println(info+` use time =`, end.Sub(now).Nanoseconds()/int64(duration), durationName)
}

func DurationToString(duration time.Duration) string {
	var durationName = ""
	if duration == time.Nanosecond {
		durationName = "ns"
	} else if duration == time.Microsecond {
		durationName = "mcs"
	} else if duration == time.Millisecond {
		durationName = "ms"
	} else if duration == time.Second {
		durationName = "s"
	} else if duration == time.Minute {
		durationName = "minute"
	} else if duration == time.Hour {
		durationName = "hour"
	}
	return durationName
}

//调用次方法 必须加上 defer 关键字！
func CountMethodTps(total float64, now time.Time, info string) {
	var end = time.Now()
	fmt.Println(info+` tps =`, total/end.Sub(now).Seconds())
}

package utils

import (
	"fmt"
	"testing"
	"time"
)

func TestGetValue(t *testing.T) {
	var v interface{}
	v = 1
	converResult := GetValue(v, nil)
	fmt.Println(converResult)
	v = int8(1)
	converResult = GetValue(v, nil)
	fmt.Println(converResult)
	v = int16(1)
	converResult = GetValue(v, nil)
	fmt.Println(converResult)
	v = int32(1)
	converResult = GetValue(v, nil)
	fmt.Println(converResult)
	v = int64(1)
	converResult = GetValue(v, nil)
	fmt.Println(converResult)
	v = "string"
	converResult = GetValue(v, nil)
	fmt.Println(converResult)
	v = time.Now()
	converResult = GetValue(v, nil)
	fmt.Println(converResult)
}

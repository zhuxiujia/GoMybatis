package GoMybatis

import (
	"testing"
	"fmt"
)

type TestResult struct {
	Name string
	Amount1  float32
	Amount2  float64
	Age1  int
	Age2  int32
	Age3  int64
	Age4  uint
	Age5  uint8
	Age6  uint16
	Age7  uint32
	Age8  uint64
	Bool bool
}

func Test_convert(t *testing.T) {
	var GoMybatisSqlResultDecoder = GoMybatisSqlResultDecoder{}
	var res = make([]map[string][]byte, 0)

	var resMap = make(map[string][]byte)
	resMap["Name"] = []byte("xiao ming")
	resMap["Amount1"] = []byte("1908.1")
	resMap["Amount2"] = []byte("1908.444")
	resMap["Age1"] = []byte("1908")
	resMap["Age2"] = []byte("1908")
	resMap["Age3"] = []byte("1908")
	resMap["Age4"] = []byte("1908")
	resMap["Age5"] = []byte("1908")
	resMap["Age6"] = []byte("1908")
	resMap["Age7"] = []byte("1908")
	resMap["Age8"] = []byte("1908")
	resMap["Bool"] = []byte("1")
	res = append(res, resMap)

	var result TestResult
	GoMybatisSqlResultDecoder.Decode(res, &result)

	fmt.Println(result)
}

func Test_Ignore_Case_Underscores(t *testing.T)  {
	var GoMybatisSqlResultDecoder = GoMybatisSqlResultDecoder{}
	var res = make([]map[string][]byte, 0)

	var resMap = make(map[string][]byte)
	resMap["name"] = []byte("xiao ming")
	resMap["Amount_1"] = []byte("1908.1")
	resMap["amount_2"] = []byte("1908.444")
	resMap["age_1"] = []byte("1908")
	resMap["age_2"] = []byte("1908")
	resMap["age_3"] = []byte("1908")
	resMap["age_4"] = []byte("1908")
	resMap["age_5"] = []byte("1908")
	resMap["age_6"] = []byte("1908")
	resMap["age_7"] = []byte("1908")
	resMap["age_8"] = []byte("1908")
	resMap["Bool"] = []byte("1")
	res = append(res, resMap)

	var result TestResult
	GoMybatisSqlResultDecoder.Decode(res, &result)

	fmt.Println(result)
}

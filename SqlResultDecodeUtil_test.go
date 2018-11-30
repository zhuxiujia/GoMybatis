package GoMybatis

import (
	"testing"
	"fmt"
	"github.com/zhuxiujia/GoMybatis/utils"
	"time"
)

type TestResult struct {
	Name    string
	Amount1 float32
	Amount2 float64
	Age1    int
	Age2    int32
	Age3    int64
	Age4    uint
	Age5    uint8
	Age6    uint16
	Age7    uint32
	Age8    uint64
	Bool    bool
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
	resMap["Bool"] = []byte("true")
	res = append(res, resMap)

	var result TestResult
	GoMybatisSqlResultDecoder.Decode(nil, res, &result)

	fmt.Println(result)
}

func Test_Ignore_Case_Underscores(t *testing.T) {
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
	GoMybatisSqlResultDecoder.Decode(nil, res, &result)

	fmt.Println(result)
}

func Test_Ignore_Case_Underscores_Tps(t *testing.T) {
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

	defer utils.CountMethodTps(100000, time.Now(), "Test_Ignore_Case_Underscores_Tps")
	for i := 0; i < 100000; i++ {
		GoMybatisSqlResultDecoder.Decode(nil, res, &result)
	}

}

func Test_Decode_Interface(t *testing.T) {

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

	var resultMap = make(map[string]ResultProperty)
	resultMap["id"] = ResultProperty{
		XMLName:  "id",
		Column:   "id",
		Property: "id",
		GoType:   "string",
	}
	resultMap["name"] = ResultProperty{
		XMLName:  "result",
		Column:   "name",
		Property: "Name",
		GoType:   "string",
	}
	resultMap["Amount_1"] = ResultProperty{
		XMLName:  "result",
		Column:   "Amount_1",
		Property: "amount_1",
		GoType:   "string",
	}
	resultMap["amount_2"] = ResultProperty{
		XMLName:  "result",
		Column:   "Amount_2",
		Property: "amount_2",
		GoType:   "string",
	}
	var result map[string]interface{}
	GoMybatisSqlResultDecoder{}.Decode(resultMap, res, &result)

	fmt.Println(result)
}

package GoMybatis

import (
	"fmt"
	"github.com/zhuxiujia/GoMybatis/utils"
	"testing"
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

//解码基本数据集
func Test_Convert_Basic_Type(t *testing.T) {
	var resMap = make(map[string][]byte)
	resMap["Amount1"] = []byte("1908")
	var resMapArray = make([]map[string][]byte, 0)
	resMapArray = append(resMapArray, resMap)

	var intResult int
	var error = GoMybatisSqlResultDecoder{}.Decode(nil, resMapArray, &intResult)
	if error != nil {
		panic(error)
	}
	fmt.Println("Test_Convert_Basic_Type,int=", intResult)

	var stringResult string
	error = GoMybatisSqlResultDecoder{}.Decode(nil, resMapArray, &stringResult)
	if error != nil {
		panic(error)
	}
	fmt.Println("Test_Convert_Basic_Type,string=", stringResult)

	var floatResult float64
	error = GoMybatisSqlResultDecoder{}.Decode(nil, resMapArray, &floatResult)
	if error != nil {
		panic(error)
	}
	fmt.Println("Test_Convert_Basic_Type,float=", floatResult)

	resMap = make(map[string][]byte)
	resMap["Date"] = []byte(time.Now().Format(time.RFC3339))
	resMapArray = make([]map[string][]byte, 0)
	resMapArray = append(resMapArray, resMap)
	var timeResult time.Time
	error = GoMybatisSqlResultDecoder{}.Decode(nil, resMapArray, &timeResult)
	if error != nil {
		panic(error)
	}
	fmt.Println("Test_Convert_Basic_Type,time=", timeResult)
}

//解码数组
func Test_Convert_Slice(t *testing.T) {
	var resMap = make(map[string][]byte)
	resMap["Amount1"] = []byte("1908")
	resMap["Amount2"] = []byte("1901")
	var resMapArray = make([]map[string][]byte, 0)
	resMapArray = append(resMapArray, resMap)

	var result []int
	var error = GoMybatisSqlResultDecoder{}.Decode(nil, resMapArray, &result)
	if error != nil {
		panic(error)
	}
	fmt.Println("Test_Convert_Basic_Type", result)
}

//解码map
func Test_Convert_Map(t *testing.T) {
	var resMap = make(map[string][]byte)
	resMap["Amount1"] = []byte("1908")
	resMap["Amount2"] = []byte("1901")
	var resMapArray = make([]map[string][]byte, 0)
	resMapArray = append(resMapArray, resMap)

	var result map[string]int
	var error = GoMybatisSqlResultDecoder{}.Decode(nil, resMapArray, &result)
	if error != nil {
		panic(error)
	}
	fmt.Println("Test_Convert_Basic_Type", result)
}

func Test_convert_struct(t *testing.T) {
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

	fmt.Println("Test_convert_struct", result)
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

	fmt.Println("Test_Ignore_Case_Underscores", result)
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

	var resultMap = make(map[string]*ResultProperty)
	resultMap["id"] = &ResultProperty{
		XMLName:  "id",
		Column:   "id",
		Property: "id",
		GoType:   "string",
	}
	resultMap["name"] = &ResultProperty{
		XMLName:  "result",
		Column:   "name",
		Property: "Name",
		GoType:   "string",
	}
	resultMap["Amount_1"] = &ResultProperty{
		XMLName:  "result",
		Column:   "Amount_1",
		Property: "amount_1",
		GoType:   "string",
	}
	resultMap["amount_2"] = &ResultProperty{
		XMLName:  "result",
		Column:   "Amount_2",
		Property: "amount_2",
		GoType:   "string",
	}
	var result map[string]interface{}
	GoMybatisSqlResultDecoder{}.Decode(resultMap, res, &result)

	fmt.Println("Test_Decode_Interface", result)
}

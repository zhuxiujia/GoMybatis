package GoMybatis

import (
	"fmt"
	"testing"
	"time"

	"github.com/zhuxiujia/GoMybatis/utils"
)

type TestResult struct {
	Name    string  `json:"name"`
	Amount1 float32 `json:"amount_1"`
	Amount2 float64 `json:"amount_2"`
	Age1    int     `json:"age_1"`
	Age2    int32   `json:"age_2"`
	Age3    int64   `json:"age_3"`
	Age4    uint    `json:"age_4"`
	Age5    uint8   `json:"age_5"`
	Age6    uint16  `json:"age_6"`
	Age7    uint32  `json:"age_7"`
	Age8    uint64  `json:"age_8"`
	Bool    bool    `json:"bool"`
}

//解码基本数据-int,string,time.Time...
func Test_Convert_Basic_Type(t *testing.T) {
	var resMap = make(map[string][]byte)
	resMap["Amount1"] = []byte("1908")
	var resMapArray = QueryResult{}
	resMapArray.append(resMap)

	var intResult int
	var error = GoMybatisSqlResultDecoder{}.Decode(nil, resMapArray, &intResult)
	if error != nil {
		t.Fatal(error)
	}
	fmt.Println("Test_Convert_Basic_Type,int=", intResult)

	var stringResult string
	error = GoMybatisSqlResultDecoder{}.Decode(nil, resMapArray, &stringResult)
	if error != nil {
		t.Fatal(error)
	}
	fmt.Println("Test_Convert_Basic_Type,string=", stringResult)

	var floatResult float64
	error = GoMybatisSqlResultDecoder{}.Decode(nil, resMapArray, &floatResult)
	if error != nil {
		t.Fatal(error)
	}
	fmt.Println("Test_Convert_Basic_Type,float=", floatResult)

	resMap = make(map[string][]byte)
	resMap["Date"] = []byte(time.Now().Format(time.RFC3339))
	resMapArray = QueryResult{}
	resMapArray.append(resMap)
	var timeResult time.Time
	error = GoMybatisSqlResultDecoder{}.Decode(nil, resMapArray, &timeResult)
	if error != nil {
		t.Fatal(error)
	}
	fmt.Println("Test_Convert_Basic_Type,time=", timeResult)
}

//解码数组
func Test_Convert_Slice(t *testing.T) {
	var resMap = make(map[string][]byte)
	resMap["Amount1"] = []byte("1908")
	resMap["Amount2"] = []byte("1901")
	var resMapArray = QueryResult{}
	resMapArray.append(resMap)

	var result []int
	var error = GoMybatisSqlResultDecoder{}.Decode(nil, resMapArray, &result)
	if error != nil {
		t.Fatal(error)
	}
	fmt.Println("Test_Convert_Slice", result)
}

//解码map
func Test_Convert_Map(t *testing.T) {
	var resMap = make(map[string][]byte)
	resMap["Amount1"] = []byte("1908")
	resMap["Amount2"] = []byte("1901")
	var resMapArray = QueryResult{}
	resMapArray.append(resMap)

	var result map[string]string
	var error = GoMybatisSqlResultDecoder{}.Decode(nil, resMapArray, &result)
	if error != nil {
		t.Fatal(error)
	}
	fmt.Println("Test_Convert_Map", result)

	resMapArray.append(resMap)

	var resultMapArray []map[string]string
	error = GoMybatisSqlResultDecoder{}.Decode(nil, resMapArray, &resultMapArray)
	if error != nil {
		t.Fatal(error)
	}
	fmt.Println("Test_Convert_Map", resultMapArray)
}

func Test_convert_struct(t *testing.T) {
	var GoMybatisSqlResultDecoder = GoMybatisSqlResultDecoder{}
	var res = QueryResult{}

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
	res.append(resMap)

	var result TestResult
	GoMybatisSqlResultDecoder.Decode(nil, res, &result)

	fmt.Println("Test_convert_struct", result)
}

func Test_Ignore_Case_Underscores(t *testing.T) {
	var GoMybatisSqlResultDecoder = GoMybatisSqlResultDecoder{}
	var res = QueryResult{}

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
	res.append(resMap)

	var result TestResult
	GoMybatisSqlResultDecoder.Decode(nil, res, &result)

	fmt.Println("Test_Ignore_Case_Underscores", result)
}

func Test_Ignore_Case_Underscores_Tps(t *testing.T) {
	var GoMybatisSqlResultDecoder = GoMybatisSqlResultDecoder{}
	var res = QueryResult{}

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
	res.append(resMap)

	var result TestResult

	defer utils.CountMethodTps(10000, time.Now(), "Test_Ignore_Case_Underscores_Tps")
	for i := 0; i < 10000; i++ {
		GoMybatisSqlResultDecoder.Decode(nil, res, &result)
	}

}

func Test_Decode_Interface(t *testing.T) {

	var res = QueryResult{}
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
	res.append(resMap)

	var resultMap = make(map[string]*ResultProperty)
	resultMap["id"] = &ResultProperty{
		XMLName:  "id",
		Column:   "id",
		Property: "id",
		LangType: "string",
	}
	resultMap["name"] = &ResultProperty{
		XMLName:  "result",
		Column:   "name",
		Property: "Name",
		LangType: "string",
	}
	resultMap["Amount_1"] = &ResultProperty{
		XMLName:  "result",
		Column:   "Amount_1",
		Property: "amount_1",
		LangType: "string",
	}
	resultMap["amount_2"] = &ResultProperty{
		XMLName:  "result",
		Column:   "Amount_2",
		Property: "amount_2",
		LangType: "string",
	}
	var result map[string]string
	GoMybatisSqlResultDecoder{}.Decode(resultMap, res, &result)

	fmt.Println("Test_Decode_Interface", result)
}

func Benchmark_Ignore_Case_Underscores(b *testing.B) {
	b.StopTimer()
	var GoMybatisSqlResultDecoder = GoMybatisSqlResultDecoder{}
	var res = QueryResult{}

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
	res.append(resMap)

	var result TestResult

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		GoMybatisSqlResultDecoder.Decode(nil, res, &result)
	}

}

//
//func TestGoMybatisSqlResultDecoder_Decode(t *testing.T) {
//	var GoMybatisSqlResultDecoder = GoMybatisSqlResultDecoder{}
//	var res = make([]map[string][]byte, 0)
//	var resMap = make(map[string][]byte)
//	resMap["name"] = []byte("xiao ming")
//	resMap["Amount_1"] = []byte("1908.1")
//	resMap["amount_2"] = []byte("1908.444")
//	resMap["age_1"] = []byte("1908")
//	resMap["age_2"] = []byte("1908")
//	resMap["age_3"] = []byte("1908")
//	resMap["age_4"] = []byte("1908")
//	resMap["age_5"] = []byte("1908")
//	resMap["age_6"] = []byte("1908")
//	resMap["age_7"] = []byte("1908")
//	resMap["age_8"] = []byte("1908")
//	resMap["Bool"] = []byte("1")
//	res = append(res, resMap)
//	var result TestResult
//	var err = GoMybatisSqlResultDecoder.Decode(nil, res, &result)
//	if err != nil {
//		t.Fatal(err)
//	}
//	if result.Name == "" ||
//		result.Amount1 == 0 ||
//		result.Amount2 == 0 ||
//		result.Age1 == 0 ||
//		result.Age2 == 0 ||
//		result.Age3 == 0 ||
//		result.Age4 == 0 ||
//		result.Age5 == 0 ||
//		result.Age6 == 0 ||
//		result.Age7 == 0 ||
//		result.Age8 == 0 ||
//		result.Bool == false {
//		t.Fatal("TestGoMybatisSqlResultDecoder_Decode fail,result not decoded!")
//	}
//	fmt.Println(result)
//}

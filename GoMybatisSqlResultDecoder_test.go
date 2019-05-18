package GoMybatis

import (
	"encoding/json"
	"github.com/zhuxiujia/GoMybatis/utils"
	"fmt"
	"testing"
	"time"
)

type TestResult struct {
	Name    string `json:"name"`
	Amount1 float32 `json:"amount_1"`
	Amount2 float64 `json:"amount_2"`
	Age1    int `json:"age_1"`
	Age2    int32 `json:"age_2"`
	Age3    int64 `json:"age_3"`
	Age4    uint `json:"age_4"`
	Age5    uint8 `json:"age_5"`
	Age6    uint16 `json:"age_6"`
	Age7    uint32 `json:"age_7"`
	Age8    uint64 `json:"age_8"`
	Bool    bool `json:"bool"`
}

//解码基本数据-int,string,time.Time...
func Test_Convert_Basic_Type(t *testing.T) {
	var resMap = make(map[string]interface{})
	resMap["Amount1"] = "1908"
	var resMapArrayData = make([]map[string]interface{}, 0)
	resMapArrayData = append(resMapArrayData, resMap)

	var b,_=json.Marshal(resMapArrayData)

	var resMapArray=string(b)

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

	resMap = make(map[string]interface{})
	resMap["Date"] = time.Now().Format(time.RFC3339)
	resMapArrayData = make([]map[string]interface{}, 0)
	resMapArrayData = append(resMapArrayData, resMap)

	b,_=json.Marshal(resMapArrayData)
	resMapArray=string(b)

	var timeResult time.Time
	error = GoMybatisSqlResultDecoder{}.Decode(nil, resMapArray, &timeResult)
	if error != nil {
		t.Fatal(error)
	}
	fmt.Println("Test_Convert_Basic_Type,time=", timeResult)
}


//解码map
func Test_Convert_Map(t *testing.T) {
	var resMap = make(map[string]interface{})
	resMap["Amount1"] = "1908"
	resMap["Amount2"] = "1901"
	var resMapArrayData = make([]map[string]interface{}, 0)

	resMapArrayData = append(resMapArrayData, resMap)
	var b,_=json.Marshal(resMapArrayData)

	var resMapArray=string(b)

	var result =map[string]string{}
	var error = GoMybatisSqlResultDecoder{}.Decode(nil, resMapArray, &result)
	if error != nil {
		t.Fatal(error)
	}
	fmt.Println("Test_Convert_Map", result)

	resMapArrayData = append(resMapArrayData, resMap)
	 b,_=json.Marshal(resMapArrayData)
	 resMapArray=string(b)


	var resultMapArray []map[string]string
	error = GoMybatisSqlResultDecoder{}.Decode(nil, resMapArray, &resultMapArray)
	if error != nil {
		t.Fatal(error)
	}
	fmt.Println("Test_Convert_Map_Array", resultMapArray)
}

func Test_convert_struct(t *testing.T) {
	var GoMybatisSqlResultDecoder = GoMybatisSqlResultDecoder{}
	var res = make([]map[string]interface{}, 0)

	var resMap = make(map[string]interface{})
	resMap["Name"] = "xiao ming"
	resMap["Amount1"] = "1908.1"
	resMap["Amount2"] = "1908.444"
	resMap["Age1"] = "1908"
	resMap["Age2"] = "1908"
	resMap["Age3"] = "1908"
	resMap["Age4"] = "1908"
	resMap["Age5"] = "1908"
	resMap["Age6"] = "1908"
	resMap["Age7"] = "1908"
	resMap["Age8"] = "1908"
	resMap["Bool"] = "true"
	res = append(res, resMap)

	var b,_=json.Marshal(res)

	var resMapArray=string(b)

	var result TestResult
	GoMybatisSqlResultDecoder.Decode(nil, resMapArray, &result)

	fmt.Println("Test_convert_struct", result)
}

func Test_Ignore_Case_Underscores(t *testing.T) {
	var GoMybatisSqlResultDecoder = GoMybatisSqlResultDecoder{}
	var res = make([]map[string]interface{}, 0)

	var resMap = make(map[string]interface{})
	resMap["name"] = "xiao ming"
	resMap["Amount_1"] = "1908.1"
	resMap["amount_2"] = "1908.444"
	resMap["age_1"] = "1908"
	resMap["age_2"] = "1908"
	resMap["age_3"] = "1908"
	resMap["age_4"] = "1908"
	resMap["age_5"] = "1908"
	resMap["age_6"] = "1908"
	resMap["age_7"] = "1908"
	resMap["age_8"] = "1908"
	resMap["Bool"] = "1"
	res = append(res, resMap)

	var b,_=json.Marshal(res)

	var resMapArray=string(b)

	var result TestResult
	GoMybatisSqlResultDecoder.Decode(nil, resMapArray, &result)

	fmt.Println("Test_Ignore_Case_Underscores", result)
}

func Test_Ignore_Case_Underscores_Tps(t *testing.T) {
	var GoMybatisSqlResultDecoder = GoMybatisSqlResultDecoder{}
	var res = make([]map[string]interface{}, 0)

	var resMap = make(map[string]interface{})
	resMap["name"] = "xiao ming"
	resMap["Amount_1"] = "1908.1"
	resMap["amount_2"] = "1908.444"
	resMap["age_1"] = "1908"
	resMap["age_2"] = "1908"
	resMap["age_3"] = "1908"
	resMap["age_4"] = "1908"
	resMap["age_5"] = "1908"
	resMap["age_6"] = "1908"
	resMap["age_7"] = "1908"
	resMap["age_8"] = "1908"
	resMap["Bool"] = "1"
	res = append(res, resMap)

	var b,_=json.Marshal(res)

	var resMapArray=string(b)

	var result TestResult

	defer utils.CountMethodTps(10000, time.Now(), "Test_Ignore_Case_Underscores_Tps")
	for i := 0; i < 10000; i++ {
		GoMybatisSqlResultDecoder.Decode(nil, resMapArray, &result)
	}

}

func Test_Decode_Interface(t *testing.T) {

	var res = make([]map[string]interface{}, 0)
	var resMap = make(map[string]interface{})
	resMap["name"] = "xiao ming"
	resMap["Amount_1"] = "1908.1"
	resMap["amount_2"] = "1908.444"
	resMap["age_1"] = "1908"
	resMap["age_2"] = "1908"
	resMap["age_3"] = "1908"
	resMap["age_4"] = "1908"
	resMap["age_5"] = "1908"
	resMap["age_6"] = "1908"
	resMap["age_7"] = "1908"
	resMap["age_8"] = "1908"
	resMap["Bool"] = "1"
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

	var b,_=json.Marshal(res)

	var resMapArray=string(b)

	var result map[string]string
	GoMybatisSqlResultDecoder{}.Decode(resultMap, resMapArray, &result)

	fmt.Println("Test_Decode_Interface", result)
}

func Benchmark_Ignore_Case_Underscores(b *testing.B) {
	b.StopTimer()
	var GoMybatisSqlResultDecoder = GoMybatisSqlResultDecoder{}
	var res = make([]map[string]interface{}, 0)

	var resMap = make(map[string]interface{})
	resMap["name"] = "xiao ming"
	resMap["Amount_1"] = "1908.1"
	resMap["amount_2"] = "1908.444"
	resMap["age_1"] = "1908"
	resMap["age_2"] = "1908"
	resMap["age_3"] = "1908"
	resMap["age_4"] = "1908"
	resMap["age_5"] = "1908"
	resMap["age_6"] = "1908"
	resMap["age_7"] = "1908"
	resMap["age_8"] = "1908"
	resMap["Bool"] = "1"
	res = append(res, resMap)


	var bt,_=json.Marshal(res)

	var resMapArray=string(bt)

	var result TestResult

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		GoMybatisSqlResultDecoder.Decode(nil, resMapArray, &result)
	}

}
package GoMybatis

import (
	"fmt"
	"github.com/zhuxiujia/GoMybatis/utils"
	"testing"
	"time"
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
	F       bool    `json:"f"`
}

//解码基本数据-int,string,time.Time...
func Test_Convert_Basic_Type(t *testing.T) {
	var resMap = make(map[string][]byte)
	resMap["Amount1"] = []byte("1908")
	var resMapArray = make([]map[string][]byte, 0)
	resMapArray = append(resMapArray, resMap)

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
	resMapArray = make([]map[string][]byte, 0)
	resMapArray = append(resMapArray, resMap)
	var timeResult time.Time
	error = GoMybatisSqlResultDecoder{}.Decode(nil, resMapArray, &timeResult)
	if error != nil {
		t.Fatal(error)
	}
	fmt.Println("Test_Convert_Basic_Type,time=", timeResult)
}

//解码map
func Test_Convert_Map(t *testing.T) {
	var resMap = make(map[string][]byte)
	resMap["Amount1"] = []byte("1908")
	resMap["Amount2"] = []byte("1901")
	var resMapArray = make([]map[string][]byte, 0)
	resMapArray = append(resMapArray, resMap)

	var result map[string]interface{}
	var error = GoMybatisSqlResultDecoder{}.Decode(nil, resMapArray, &result)
	if error != nil {
		t.Fatal(error)
	}
	fmt.Println("Test_Convert_Map", result)

	resMapArray = append(resMapArray, resMap)

	var resultMapArray []map[string]interface{}
	error = GoMybatisSqlResultDecoder{}.Decode(nil, resMapArray, &resultMapArray)
	if error != nil {
		t.Fatal(error)
	}
	fmt.Println("Test_Convert_Map", resultMapArray)
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
	resMap["Age5"] = []byte("1")
	resMap["Age6"] = []byte("1908")
	resMap["Age7"] = []byte("1908")
	resMap["Age8"] = []byte("1908")
	resMap["Bool"] = []byte("true")
	res = append(res, resMap)

	var result TestResult
	GoMybatisSqlResultDecoder.Decode(nil, res, &result)

	if result.Name != string(resMap["Name"]) {
		panic("convert_struct Name fail")
	}
	if result.Amount1 != 1908.1 {
		panic("convert_struct Amount1 fail")
	}
	if result.Amount2 != 1908.444 {
		panic("convert_struct Amount2 fail")
	}

	if result.Age1 != 1908 {
		panic("convert_struct Age1 fail")
	}
	if result.Age2 != 1908 {
		panic("convert_struct Age1 fail")
	}
	if result.Age3 != 1908 {
		panic("convert_struct Age1 fail")
	}
	if result.Age4 != 1908 {
		panic("convert_struct Age1 fail")
	}
	if result.Age5 != 1 {
		panic("convert_struct Age1 fail")
	}
	if result.Age6 != 1908 {
		panic("convert_struct Age1 fail")
	}
	if result.Age7 != 1908 {
		panic("convert_struct Age1 fail")
	}
	if result.Age8 != 1908 {
		panic("convert_struct Age1 fail")
	}
	if result.Bool != true {
		panic("convert_struct Bool fail")
	}

	fmt.Println("Test_convert_struct", result)
}

func Test_convert_html_struct(t *testing.T) {
	var GoMybatisSqlResultDecoder = GoMybatisSqlResultDecoder{}
	var res = make([]map[string][]byte, 0)

	var resMap = make(map[string][]byte)
	resMap["Name"] = []byte("<p>adfd</p><section class=\"\n\twwei - editor \"><fieldset style=\"\n\tclear: both;padding: 0 px;border: 0 px none;margin: 1e m 0 px 0.5e m;\n\t\"><section style=\"\n\tborder - top: 2 px solid rgb(134, 110, 187);font - size: 1e m;font - weight: inherit;text - decoration: inherit;color: rgb(255, 255, 255);box - sizing: border - box;border - bottom - color: rgb(134, 110, 187);border - left - color: rgb(134, 110, 187);border - right - color: rgb(134, 110, 187);\n\t\"><section class=\"\n\twweibrush \" data-brushtype=\"\n\ttext \" style=\"\n\tpadding: 0 px 0.5e m;display: inline - block;font - size: 16 px;background - color: rgb(134, 110, 187);color: rgb(255, 255, 255);\n\t\">微信编辑器</section></section></fieldset></section><p><br/></p><section class=\"\n\twwei - editor \"><section style=\"\n\tfont - size: 1e m;white - space: normal;margin: 1e m 0 px 0.8e m;text - align: center;vertical - align: middle;\n\t\"><section style=\"\n\theight: 0 px;margin: 0 px 1e m;border - width: 1.5e m;border - style: solid;border - image: initial;border - left - color: transparent!important;border - right - color: transparent!important;border - bottom - color: rgb(134, 110, 187);border - top - color: rgb(134, 110, 187);\n\t\"></section><section style=\"\n\theight: 0 px;margin: -2.75e m 1.65e m;border - width: 1.3e m;border - style: solid;border - color: rgb(255, 255, 255) transparent;\n\t\"></section><section style=\"\n\theight: 0 px;margin: 0.45e m 2.1e m;vertical - align: middle;border - width: 1.1e m;border - style: solid;border - image: initial;border - left - color: transparent!important;border - right - color: transparent!important;border - bottom - color: rgb(134, 110, 187);border - top - color: rgb(134, 110, 187);\n\t\"><section class=\"\n\twweibrush \" data-brushtype=\"\n\ttext \" placeholder=\"\n\t一行短标题 \" style=\"\n\tmax - height: 1e m;padding: 0 px;margin - top: -0.5e m;color: rgb(255, 255, 255);font - size: 1.2e m;line - height: 1e m;overflow: hidden;\n\t\">一行短标题</section></section></section></section><p><br/></p><p><br/></p>")
	resMap["Amount1"] = []byte("1908.1")
	resMap["Amount2"] = []byte("1908.444")
	resMap["Age1"] = []byte("1908")
	resMap["Age2"] = []byte("1908")
	resMap["Age3"] = []byte("1908")
	resMap["Age4"] = []byte("1908")
	resMap["Age5"] = []byte("1")
	resMap["Age6"] = []byte("1908")
	resMap["Age7"] = []byte("1908")
	resMap["Age8"] = []byte("1908")
	resMap["Bool"] = []byte("true")
	res = append(res, resMap)

	var result TestResult
	GoMybatisSqlResultDecoder.Decode(nil, res, &result)

	//if result.Name != string(resMap["Name"]) {
	//	panic("convert_struct Name fail")
	//}
	if result.Amount1 != 1908.1 {
		panic("convert_struct Amount1 fail")
	}
	if result.Amount2 != 1908.444 {
		panic("convert_struct Amount2 fail")
	}

	if result.Age1 != 1908 {
		panic("convert_struct Age1 fail")
	}
	if result.Age2 != 1908 {
		panic("convert_struct Age1 fail")
	}
	if result.Age3 != 1908 {
		panic("convert_struct Age1 fail")
	}
	if result.Age4 != 1908 {
		panic("convert_struct Age1 fail")
	}
	if result.Age5 != 1 {
		panic("convert_struct Age1 fail")
	}
	if result.Age6 != 1908 {
		panic("convert_struct Age1 fail")
	}
	if result.Age7 != 1908 {
		panic("convert_struct Age1 fail")
	}
	if result.Age8 != 1908 {
		panic("convert_struct Age1 fail")
	}
	if result.Bool != true {
		panic("convert_struct Bool fail")
	}

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

	defer utils.CountMethodTps(10000, time.Now(), "Test_Ignore_Case_Underscores_Tps")
	for i := 0; i < 10000; i++ {
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
		LangType: "string",
	}
	resultMap["name"] = &ResultProperty{
		XMLName:  "result",
		Column:   "name",
		LangType: "string",
	}
	resultMap["Amount_1"] = &ResultProperty{
		XMLName:  "result",
		Column:   "Amount_1",
		LangType: "string",
	}
	resultMap["amount_2"] = &ResultProperty{
		XMLName:  "result",
		Column:   "Amount_2",
		LangType: "string",
	}
	var result map[string]string
	GoMybatisSqlResultDecoder{}.Decode(resultMap, res, &result)

	fmt.Println("Test_Decode_Interface", result)
}

func Benchmark_Ignore_Case_Underscores(b *testing.B) {
	b.StopTimer()
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

package GoMybatis

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func Test_Load_Xml(t *testing.T) {
	//读取mapper xml文件
	file, err := os.Open("example/Example_ActivityMapper.xml")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	bytes, _ := ioutil.ReadAll(file)
	var xmlItems = LoadMapperXml(bytes)
	if xmlItems==nil{
		t.Fatal(`Test_Load_Xml fail,LoadMapperXml "example/Example_ActivityMapper.xml"`)
	}
	fmt.Println(xmlItems)
}

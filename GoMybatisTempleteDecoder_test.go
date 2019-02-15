package GoMybatis

import (
	"fmt"
	"github.com/zhuxiujia/GoMybatis/example"
	"io/ioutil"
	"testing"
	"time"
)

type ExampleActivityMapper struct {
	SelectTemplete func(name string, session Session) ([]example.Activity, error) `mapperParams:"name,session"`
	InsertTemplete func(arg example.Activity, session Session) (int64, error)     `mapperParams:"arg,session"`
	UpdateTemplete func(arg example.Activity, session Session) (int64, error)     `mapperParams:"name,session"`
	DeleteTemplete func(name string, session Session) (int64, error)              `mapperParams:"name,session"`
}

//初始化mapper文件和结构体
var exampleActivityMapper = ExampleActivityMapper{}

func init() {
	bytes, err := ioutil.ReadFile("example/Example_ActivityMapper.xml")
	if err != nil {
		panic(err)
	}
	var xmlItems = LoadMapperXml(bytes)
	if xmlItems == nil {
		panic(`Test_Load_Xml fail,LoadMapperXml "example/Example_ActivityMapper.xml"`)
	}

	var decoder = GoMybatisTempleteDecoder{}
	err = decoder.DecodeTree(xmlItems, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(xmlItems)

	var engine = GoMybatisEngine{}.New()
	//mysql链接格式为         用户名:密码@(数据库链接地址:端口)/数据库名称   例如root:123456@(***.mysql.rds.aliyuncs.com:3306)/test
	err = engine.Open("mysql", "") //此处请按格式填写你的mysql链接，这里用*号代替
	if err != nil {
		panic(err)
	}
	engine.SetLogEnable(true)
	engine.WriteMapperPtr(&exampleActivityMapper, bytes)
}

type TempleteSession struct {
	Session
}

func (it *TempleteSession) Id() string {
	return "sadf"
}
func (it *TempleteSession) Query(sqlorArgs string) ([]map[string][]byte, error) {
	resultsSlice := make([]map[string][]byte, 0)

	result := make(map[string][]byte)
	result["name"] = []byte("活动1")
	result["id"] = []byte("125")
	result["pc_link"] = []byte("http://www.baidu.com")
	result["h5_link"] = []byte("http://www.baidu.com")
	result["remark"] = []byte("活动1")
	resultsSlice = append(resultsSlice, result)
	return resultsSlice, nil
}
func (it *TempleteSession) Exec(sqlorArgs string) (*Result, error) {
	var result = Result{
		LastInsertId: 1,
		RowsAffected: 1,
	}
	return &result, nil
}
func (it *TempleteSession) Rollback() error {
	return nil
}
func (it *TempleteSession) Commit() error {
	return nil
}
func (it *TempleteSession) Begin() error {
	return nil
}
func (it *TempleteSession) Close() {

}

func TestGoMybatisTempleteDecoder_Decode(t *testing.T) {
	var decoder = GoMybatisTempleteDecoder{}
	var mapElements = make([]ElementItem, 0)
	mapElements = append(mapElements, ElementItem{})
	var BaseResultMap = MapperXml{
		Tag: "resultMap",
		Id:  "BaseResultMap",
		ElementItems: []ElementItem{
			{
				ElementType: "id",
				Propertys: map[string]string{
					"column":   "id",
					"property": "id",
				},
			},
			{
				ElementType: "result",
				Propertys: map[string]string{
					"column":   "name",
					"property": "name",
					"langType": "string",
				},
			},
			{
				ElementType: "result",
				Propertys: map[string]string{
					"column":   "pc_link",
					"property": "pcLink",
					"langType": "string",
				},
			},
			{
				ElementType: "result",
				Propertys: map[string]string{
					"column":   "h5_link",
					"property": "h5Link",
					"langType": "string",
				},
			},
			{
				ElementType: "result",
				Propertys: map[string]string{
					"column":   "remark",
					"property": "remark",
					"langType": "string",
				},
			},
			{
				ElementType: "result",
				Propertys: map[string]string{
					"column":         "version",
					"property":       "version",
					"langType":       "int",
					"enable_version": "true",
				},
			},
			{
				ElementType: "result",
				Propertys: map[string]string{
					"column":   "create_time",
					"property": "createTime",
					"langType": "time.Time",
				},
			},
			{
				ElementType: "result",
				Propertys: map[string]string{
					"column":               "delete_flag",
					"property":             "deleteFlag",
					"langType":             "int",
					"enable_version":       "true",
					"enable_logic_delete":  "true",
					"logic_deleted_value":  "1",
					"logic_undelete_value": "0",
				},
			},
		},
	}
	fmt.Println(BaseResultMap)

	var xml = MapperXml{
		Tag: "selectTemplete",
		Propertys: map[string]string{
			"table":   "biz_activity",
			"columns": "*",
			"wheres":  "name?name = #{name}",
		},
		ElementItems: []ElementItem{},
	}

	var e = decoder.DecodeTree(map[string]*MapperXml{"m": &xml}, nil)
	if e != nil {
		t.Fatal(e)
	}
	fmt.Println(xml.ElementItems)
}

type El struct {
	Els []El

	String string

	test string

	prefix          string
	suffix          string
	suffixOverrides string

	separator  string
	collection string
	item       string
	index      string
	open       string
}

func Test_create_conf(t *testing.T) {
	var els = []El{
		{
			test:   "name != null",
			String: "and name like #{pattern}",
		},
		{
			test:   "startTime != null",
			String: "and startTime = #{startTime}",
		},
		{
			String: "order by desc",
		},
		{
			test:   "page != null and size != null",
			String: "limit #{page}, #{size}",
		},
		{
			prefix:          "(",
			suffix:          ")",
			suffixOverrides: ",",
			Els: []El{
				{
					test:   "page != null and size != null",
					String: "limit #{page}, #{size}",
				},
			},
		},
		{
			open:       "(",
			collection: "ids",
		},
	}
	fmt.Println(els[0].test)
}

func TestGoMybatisTempleteDecoder_Create(t *testing.T) {
	var act = example.Activity{
		Id:   "123",
		Name: "test",
	}
	var session = TempleteSession{}
	n, err := exampleActivityMapper.InsertTemplete(act, &session)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("updateNum", n)
	time.Sleep(time.Second)
}

func TestGoMybatisTempleteDecoder_Select(t *testing.T) {
	var session = TempleteSession{}
	n, err := exampleActivityMapper.SelectTemplete("test", &session)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("updateNum", n)
	time.Sleep(time.Second)
}

func TestGoMybatisTempleteDecoder_Update(t *testing.T) {
	var act = example.Activity{
		Id:   "123",
		Name: "test",
	}
	var session = TempleteSession{}
	n, err := exampleActivityMapper.UpdateTemplete(act, &session)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("updateNum", n)
	time.Sleep(time.Second)
}

func TestGoMybatisTempleteDecoder_Delete(t *testing.T) {
	var session = TempleteSession{}
	n, err := exampleActivityMapper.DeleteTemplete("test", &session)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("updateNum", n)
	time.Sleep(time.Second)
}

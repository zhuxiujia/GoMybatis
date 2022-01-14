package GoMybatis

import (
	"fmt"
	"github.com/zhuxiujia/GoMybatis/example"
	"github.com/zhuxiujia/GoMybatis/stmt"
	"github.com/zhuxiujia/GoMybatis/tx"
	"testing"
	"time"
)

type ExampleActivityMapper struct {
	SelectTemplate func(name string, session Session) ([]example.Activity, error) `args:"name,session"`
	InsertTemplate func(args []example.Activity, session Session) (int64, error)  `args:"args,session"`
	UpdateTemplate func(arg example.Activity, session Session) (int64, error)     `args:"name,session"`
	DeleteTemplate func(name string, session Session) (int64, error)              `args:"name,session"`
}

//初始化mapper文件和结构体
var exampleActivityMapper = ExampleActivityMapper{}

func getMapper() ExampleActivityMapper {
	initMapper()
	return exampleActivityMapper
}

func initMapper() {
	bytes := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE mapper PUBLIC "-//mybatis.org//DTD Mapper 3.0//EN"
        "https://raw.githubusercontent.com/zhuxiujia/GoMybatis/master/mybatis-3-mapper.dtd">
<mapper>
    <!--logic_delete_key 逻辑删除字段-->
    <!--logic_deleted_value 逻辑删除已删除字段-->
    <!--logic_undelete_value 逻辑删除 未删除字段-->
    <!--version_key 乐观锁版本字段,支持int,int8,int16,int32,int64-->
    <resultMap id="BaseResultMap">
        <id column="id" langType="id"/>
        <result column="name" property="name" langType="string"/>
        <result column="pc_link" property="pcLink" langType="string"/>
        <result column="h5_link" property="h5Link" langType="string"/>
        <result column="remark" property="remark" langType="string"/>
        <result column="version" property="version" langType="int" version_enable="true"/>
        <result column="create_time" property="createTime" langType="time.Time"/>
        <result column="delete_flag" property="deleteFlag" langType="int" logic_enable="true" logic_undelete="1" logic_deleted="0"/>
    </resultMap>
    <!--模板标签: columns wheres sets 支持逗号','分隔表达式，name?name = #{name}为判空表达式-->
    <!--插入模板:默认id="insertTemplate,test="field != null",where自动设置逻辑删除字段" -->
    <!--查询模板:默认id="selectTemplate,where自动设置逻辑删除字段-->
    <!--更新模板:默认id="updateTemplate,set自动设置乐观锁版本号-->
    <!--删除模板:默认id="deleteTemplate,where自动设置逻辑删除字段-->
    <insertTemplate tables="biz_activity" />
    <selectTemplate tables="biz_activity" wheres="name?name = #{name}" columns=""/>
    <updateTemplate tables="biz_activity" sets="name?name = #{name}" wheres="name?name = #{name}"/>
    <deleteTemplate tables="biz_activity" wheres="name?name = #{name}"/>
  </mapper>
`)

	var err error

	var xmlItems = LoadMapperXml(bytes)
	if xmlItems == nil {
		panic(`Test_Load_Xml fail,LoadMapperXml "example/Example_ActivityMapper.xml"`)
	}

	var decoder = GoMybatisTemplateDecoder{}
	err = decoder.DecodeTree(xmlItems, nil)
	if err != nil {
		panic(err)
	}

	var engine = GoMybatisEngine{}.New()
	//mysql链接格式为         用户名:密码@(数据库链接地址:端口)/数据库名称   例如root:123456@(***.mysql.rds.aliyuncs.com:3306)/test
	_, err = engine.Open("mysql", "root:123456@(localhost:3306)/test") //此处请按格式填写你的mysql链接，这里用*号代替
	if err != nil {
		panic(err)
	}
	engine.SetLogEnable(false)
	engine.WriteMapperPtr(&exampleActivityMapper, bytes)
}

type TemplateSession struct {
	Session
}

func (it *TemplateSession) Id() string {
	return "sadf"
}
func (it *TemplateSession) Query(sqlorArgs string) ([]map[string][]byte, error) {
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
func (it *TemplateSession) Exec(sqlorArgs string) (*Result, error) {
	var result = Result{
		LastInsertId: 1,
		RowsAffected: 1,
	}
	return &result, nil
}
func (it *TemplateSession) QueryPrepare(sqlorArgs string, args ...interface{}) ([]map[string][]byte, error) {
	return nil, nil
}
func (it *TemplateSession) ExecPrepare(sqlorArgs string, args ...interface{}) (*Result, error) {
	var result = Result{
		LastInsertId: 1,
		RowsAffected: 1,
	}
	return &result, nil
}

func (it *TemplateSession) Rollback() error {
	return nil
}
func (it *TemplateSession) Commit() error {
	return nil
}
func (it *TemplateSession) Begin(p *tx.Propagation) error {
	return nil
}
func (it *TemplateSession) Close() {

}
func (it *TemplateSession) StmtConvert() (stmt.StmtIndexConvert, error) {
	return stmt.BuildStmtConvert("mysql")
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

func TestGoMybatisTemplateDecoder_Create(t *testing.T) {
	var act = example.Activity{
		Id:         "123",
		Uuid:       "uu",
		Name:       "test",
		PcLink:     "pc",
		H5Link:     "h5",
		Remark:     "remark",
		Version:    0,
		CreateTime: time.Now(),
		DeleteFlag: 1,
	}
	var args = []example.Activity{
		act,
		act,
		act,
		act,
		act,
	}
	var session = TemplateSession{}
	n, err := getMapper().InsertTemplate(args, &session)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("updateNum", n)
	time.Sleep(time.Second)
}

func TestGoMybatisTemplateDecoder_Select(t *testing.T) {
	var session = TemplateSession{}
	n, err := getMapper().SelectTemplate("test", &session)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("updateNum", n)
	time.Sleep(time.Second)
}

func TestGoMybatisTemplateDecoder_Update(t *testing.T) {
	var act = example.Activity{
		Id:   "123",
		Name: "test",
	}
	var session = TemplateSession{}
	n, err := getMapper().UpdateTemplate(act, &session)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("updateNum", n)
	time.Sleep(time.Second)
}

func TestGoMybatisTemplateDecoder_Delete(t *testing.T) {
	var session = TemplateSession{}
	n, err := getMapper().DeleteTemplate("test", &session)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("updateNum", n)
	time.Sleep(time.Second)
}

func TestInit(t *testing.T) {
	initMapperTest()
}

func initMapperTest() {
	bytes := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE mapper PUBLIC "-//mybatis.org//DTD Mapper 3.0//EN"
        "https://raw.githubusercontent.com/zhuxiujia/GoMybatis/master/mybatis-3-mapper.dtd">
<mapper>
    <!--logic_delete_key 逻辑删除字段-->
    <!--logic_deleted_value 逻辑删除已删除字段-->
    <!--logic_undelete_value 逻辑删除 未删除字段-->
    <!--version_key 乐观锁版本字段,支持int,int8,int16,int32,int64-->
    <resultMap id="BaseResultMap">
        <id column="id" langType="id"/>
        <result column="name" property="name" langType="string"/>
        <result column="pc_link" property="pcLink" langType="string"/>
        <result column="h5_link" property="h5Link" langType="string"/>
        <result column="remark" property="remark" langType="string"/>
        <result column="version" property="version" langType="int" version_enable="true"/>
        <result column="create_time" property="createTime" langType="time.Time"/>
        
    </resultMap>
    <!--模板标签: columns wheres sets 支持逗号','分隔表达式，name?name = #{name}为判空表达式-->
    <!--插入模板:默认id="insertTemplate,test="field != null",where自动设置逻辑删除字段" -->
    <!--查询模板:默认id="selectTemplate,where自动设置逻辑删除字段-->
    <!--更新模板:默认id="updateTemplate,set自动设置乐观锁版本号-->
    <!--删除模板:默认id="deleteTemplate,where自动设置逻辑删除字段-->
    <insertTemplate tables="biz_activity" />
    <selectTemplate tables="biz_activity" wheres="name?name = #{name}" columns=""/>
    <updateTemplate tables="biz_activity" sets="name?name = #{name}" wheres="name?name = #{name}"/>
    <deleteTemplate tables="biz_activity" wheres="name?name = #{name}"/>
  </mapper>
`)

	var err error

	var xmlItems = LoadMapperXml(bytes)
	if xmlItems == nil {
		panic(`Test_Load_Xml fail,LoadMapperXml "example/Example_ActivityMapper.xml"`)
	}

	var decoder = GoMybatisTemplateDecoder{}
	err = decoder.DecodeTree(xmlItems, nil)
	if err != nil {
		panic(err)
	}

	var engine = GoMybatisEngine{}.New()
	//mysql链接格式为         用户名:密码@(数据库链接地址:端口)/数据库名称   例如root:123456@(***.mysql.rds.aliyuncs.com:3306)/test
	_, err = engine.Open("mysql", "root:123456@(localhost:3306)/test") //此处请按格式填写你的mysql链接，这里用*号代替
	if err != nil {
		panic(err)
	}
	engine.SetLogEnable(false)
	engine.WriteMapperPtr(&exampleActivityMapper, bytes)
}

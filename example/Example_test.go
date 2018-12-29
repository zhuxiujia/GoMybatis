package example

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/zhuxiujia/GoMybatis"
	"github.com/zhuxiujia/GoMybatis/utils"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

//定义mapper文件的接口和结构体
// 支持基本类型(int,string,time.Time,float...且需要指定参数名称`mapperParams:"name"以逗号隔开，且位置要和实际参数相同)
//参数中包含有*GoMybatis.Session的类型，用于自定义事务
//自定义结构体参数（属性必须大写）
//返回中必须有error
// 函数return必须为error 为返回错误信息
type ExampleActivityMapper struct {
	SelectByIds       func(ids []string) ([]Activity, error) `mapperParams:"ids"`
	SelectAll         func() ([]map[string]string, error)
	SelectByCondition func(name string, startTime time.Time, endTime time.Time, page int, size int) ([]Activity, error) `mapperParams:"name,startTime,endTime,page,size"`
	UpdateById        func(session *GoMybatis.Session, arg Activity) (int64, error)
	Insert            func(arg Activity) (int64, error)
	CountByCondition  func(name string, startTime time.Time, endTime time.Time) (int, error) `mapperParams:"name,startTime,endTime"`
	DeleteById        func(id string) (int64, error)                                         `mapperParams:"id"`
	Choose            func(deleteFlag int) ([]Activity, error)                               `mapperParams:"deleteFlag"`
	SelectLinks       func(column string) ([]Activity, error)                                `mapperParams:"column"`
}

//初始化mapper文件和结构体
func InitMapperByLocalSession() *ExampleActivityMapper {
	var err error
	//mysql链接格式为         用户名:密码@(数据库链接地址:端口)/数据库名称   例如root:123456@(***.mysql.rds.aliyuncs.com:3306)/test
	engine, err := GoMybatis.Open("mysql", MysqlUri) //此处请按格式填写你的mysql链接，这里用*号代替
	if err != nil {
		panic(err.Error())
	}
	//读取mapper xml文件
	file, err := os.Open("Example_ActivityMapper.xml")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	bytes, _ := ioutil.ReadAll(file)
	var exampleActivityMapper ExampleActivityMapper
	//设置对应的mapper xml文件
	GoMybatis.WriteMapperPtrByEngine(&exampleActivityMapper, bytes, engine, true)
	return &exampleActivityMapper
}

//插入
func Test_inset(t *testing.T) {
	if MysqlUri == "" || MysqlUri == "*" {
		fmt.Println("no database url define in MysqlConfig.go , you must set the mysql link!")
		return
	}
	//初始化mapper文件
	var exampleActivityMapper = InitMapperByLocalSession()
	//使用mapper
	var result, err = exampleActivityMapper.Insert(Activity{Id: "171", Name: "test_insret", CreateTime: time.Now(), DeleteFlag: 1})
	if err != nil {
		panic(err)
	}
	fmt.Println("result=", result)
}

//修改
//本地事务使用例子
func Test_update(t *testing.T) {
	if MysqlUri == "" || MysqlUri == "*" {
		fmt.Println("no database url define in MysqlConfig.go , you must set the mysql link!")
		return
	}
	//初始化mapper文件
	exampleActivityMapperImpl := InitMapperByLocalSession()
	var activityBean = Activity{
		Id:   "171",
		Name: "rs168",
	}
	var updateNum, e = exampleActivityMapperImpl.UpdateById(nil, activityBean) //sessionId 有值则使用已经创建的session，否则新建一个session
	fmt.Println("updateNum=", updateNum)
	if e != nil {
		panic(e)
	}
}

//删除
func Test_delete(t *testing.T) {
	if MysqlUri == "" || MysqlUri == "*" {
		fmt.Println("no database url define in MysqlConfig.go , you must set the mysql link!")
		return
	}
	//初始化mapper文件
	var exampleActivityMapperImpl = InitMapperByLocalSession()
	//使用mapper
	var result, err = exampleActivityMapperImpl.DeleteById("171")
	if err != nil {
		panic(err)
	}
	fmt.Println("result=", result)
}

//查询
func Test_select(t *testing.T) {
	if MysqlUri == "" || MysqlUri == "*" {
		fmt.Println("no database url define in MysqlConfig.go , you must set the mysql link!")
		return
	}
	//初始化mapper文件
	var exampleActivityMapperImpl = InitMapperByLocalSession()
	//使用mapper
	var result, err = exampleActivityMapperImpl.SelectByCondition("注册", time.Time{}, time.Time{}, 0, 2000)
	if err != nil {
		panic(err)
	}
	fmt.Println("result=", result)
}

//查询
func Test_select_all(t *testing.T) {
	if MysqlUri == "" || MysqlUri == "*" {
		fmt.Println("no database url define in MysqlConfig.go , you must set the mysql link!")
		return
	}
	//初始化mapper文件
	var exampleActivityMapperImpl = InitMapperByLocalSession()
	//使用mapper
	var result, err = exampleActivityMapperImpl.SelectAll()
	if err != nil {
		panic(err)
	}
	fmt.Println("result=", result)
}

//查询
func Test_count(t *testing.T) {
	if MysqlUri == "" || MysqlUri == "*" {
		fmt.Println("no database url define in MysqlConfig.go , you must set the mysql link!")
		return
	}
	//初始化mapper文件
	var exampleActivityMapperImpl = InitMapperByLocalSession()
	//使用mapper
	var result, err = exampleActivityMapperImpl.CountByCondition("", time.Time{}, time.Time{})
	if err != nil {
		panic(err)
	}
	fmt.Println("result=", result)
}

//本地GoMybatis使用例子
func Test_ForEach(t *testing.T) {
	if MysqlUri == "" || MysqlUri == "*" {
		fmt.Println("no database url define in MysqlConfig.go , you must set the mysql link!")
		return
	}
	//初始化mapper文件
	var exampleActivityMapperImpl = InitMapperByLocalSession()
	//使用mapper
	var ids = []string{"1", "2"}
	var result, err = exampleActivityMapperImpl.SelectByIds(ids)
	if err != nil {
		panic(err)
	}
	fmt.Println("result=", result)
}

//本地事务使用例子
func Test_local_Transation(t *testing.T) {
	if MysqlUri == "" || MysqlUri == "*" {
		fmt.Println("no database url define in MysqlConfig.go , you must set the mysql link!")
		return
	}
	//初始化mapper文件
	exampleActivityMapperImpl := InitMapperByLocalSession()
	//使用事务
	var session = GoMybatis.DefaultSessionFactory.NewSession(GoMybatis.SessionType_Default, nil)
	session.Begin() //开启事务
	var activityBean = Activity{
		Id:   "170",
		Name: "rs168-8",
	}
	var updateNum, e = exampleActivityMapperImpl.UpdateById(&session, activityBean) //sessionId 有值则使用已经创建的session，否则新建一个session
	fmt.Println("updateNum=", updateNum)
	if e != nil {
		panic(e)
	}
	session.Commit() //提交事务
	session.Close()  //关闭事务
}

//远程事务示例，可用于分布式微服务(单数据库，多个微服务)
func Test_Remote_Transation(t *testing.T) {
	if MysqlUri == "" || MysqlUri == "*" {
		fmt.Println("no database url define in MysqlConfig.go , you must set the mysql link!")
		return
	}
	//启动GoMybatis独立节点事务服务器，通过rpc调用
	var remoteAddr = "127.0.0.1:17235"
	go GoMybatis.ServerTransationTcp(remoteAddr, MysqlDriverName, MysqlUri)

	//开始使用
	//初始化mapper文件
	var exampleActivityMapperImpl = InitMapperByLocalSession()

	//关键，使用远程Session替换本地Session调用
	var transationRMSession = GoMybatis.DefaultSessionFactory.NewSession(GoMybatis.SessionType_TransationRM, &GoMybatis.TransationRMClientConfig{
		Addr:          remoteAddr,
		RetryTime:     3,
		TransactionId: utils.CreateUUID(),
		Status:        GoMybatis.Transaction_Status_NO,
	})

	//开启远程事务
	transationRMSession.Begin()
	//使用mapper
	var activityBean = Activity{
		Id:   "170",
		Name: "rs168-11",
	}
	var _, err = exampleActivityMapperImpl.UpdateById(&transationRMSession, activityBean)
	if err != nil {
		panic(err)
	}
	//提交远程事务
	transationRMSession.Commit()
	//回滚远程事务
	//transationRMSession.Rollback()

	transationRMSession.Close()
}

func Test_choose(t *testing.T) {
	if MysqlUri == "" || MysqlUri == "*" {
		fmt.Println("no database url define in MysqlConfig.go , you must set the mysql link!")
		return
	}
	//初始化mapper文件
	var exampleActivityMapperImpl = InitMapperByLocalSession()
	//使用mapper
	var result, err = exampleActivityMapperImpl.Choose(1)
	if err != nil {
		panic(err)
	}
	fmt.Println("result=", result)
}

//查询
func Test_include_sql(t *testing.T) {
	if MysqlUri == "" || MysqlUri == "*" {
		fmt.Println("no database url define in MysqlConfig.go , you must set the mysql link!")
		return
	}
	//初始化mapper文件
	var exampleActivityMapperImpl = InitMapperByLocalSession()
	//使用mapper
	var result, err = exampleActivityMapperImpl.SelectLinks("name")
	if err != nil {
		panic(err)
	}
	fmt.Println("result=", result)
}

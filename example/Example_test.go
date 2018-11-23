package example

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/zhuxiujia/GoMybatis"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

//定义mapper文件的接口和结构体
// 支持基本类型(int,string,time.Time,float...且需要指定参数名称`mapperParams:"name"以逗号隔开，且位置要和实际参数相同)
//自定义结构体参数（属性必须大写）
//参数中除了session指针外，为指针类型的皆为返回数据
// 函数return必须为error 为返回错误信息
type ExampleActivityMapperImpl struct {
	SelectByIds       func(ids []string,result *[]Activity) error `mapperParams:"ids"`
	SelectAll         func(result *[]Activity) error
	SelectByCondition func(name string, startTime time.Time, endTime time.Time, page int, size int, result *[]Activity) error `mapperParams:"name,startTime,endTime,page,size"`
	UpdateById        func(session *GoMybatis.Session, arg Activity, result *int64) error                                     //只要参数中包含有*GoMybatis.Session的类型，框架默认使用传入的session对象，用于自定义事务
	Insert            func(arg Activity, result *int64) error
	CountByCondition  func(name string, startTime time.Time, endTime time.Time, result *int) error `mapperParams:"name,startTime,endTime"`
}

//初始化mapper文件和结构体
func InitMapperByLocalSession() ExampleActivityMapperImpl {
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
	var exampleActivityMapperImpl ExampleActivityMapperImpl
	//设置对应的mapper xml文件
	GoMybatis.UseProxyMapperByEngine(&exampleActivityMapperImpl, bytes, engine)
	return exampleActivityMapperImpl
}

//本地GoMybatis使用例子
func Test_main(t *testing.T) {
	//初始化mapper文件
	var exampleActivityMapperImpl = InitMapperByLocalSession()
	//使用mapper
	var result []Activity
	var err = exampleActivityMapperImpl.SelectByCondition("", time.Time{}, time.Time{}, 0, 2000, &result)
	if err != nil {
		panic(err)
	}
	fmt.Println("result=", result)
}

//本地GoMybatis使用例子
func Test_ForEach(t *testing.T) {
	//初始化mapper文件
	var exampleActivityMapperImpl = InitMapperByLocalSession()
	//使用mapper
	var result []Activity

	var ids=[]string{"170"}
	var err = exampleActivityMapperImpl.SelectByIds(ids, &result)
	if err != nil {
		panic(err)
	}
	fmt.Println("result=", result)
}

//本地事务使用例子
func Test_local_Transation(t *testing.T) {
	//初始化mapper文件
	exampleActivityMapperImpl := InitMapperByLocalSession()
	//使用事务
	session := *GoMybatis.DefaultSessionFactory.NewSession()
	session.Begin() //开启事务
	var activityBean = Activity{
		Id:   "170",
		Name: "rs168-8",
	}
	var updateNum int64 = 0
	var e = exampleActivityMapperImpl.UpdateById(&session, activityBean, &updateNum) //sessionId 有值则使用已经创建的session，否则新建一个session
	fmt.Println("updateNum=", updateNum)
	if e != nil {
		fmt.Println(e)
	}
	session.Commit() //提交事务
	session.Close()  //关闭事务
}

//远程事务示例，可用于分布式微服务(单数据库，多个微服务)
func Test_Remote_Transation(t *testing.T) {

	//启动GoMybatis独立节点事务服务器，通过msgpackrpc调用(msgpack是一种多语言序列化格式，类似json但是速度快于go默认的gob编码)
	var addr = "127.0.0.1:17235"
	go GoMybatis.ServerTcp(addr, MysqlDriverName, MysqlUri)

	var TransationRMClient = GoMybatis.TransationRMClient{
		RetryTime: 3,
		Addr:      addr,
	}

	//这里是关键，使用原程的transationRMSession（即Transation Resource Manager）替换LocalSession本地的session调用
	var transationRMSession = *GoMybatis.TransationRMSession{}.New("", &TransationRMClient, GoMybatis.Transaction_Status_NO)

	//初始化mapper文件
	var exampleActivityMapperImpl = InitMapperByLocalSession()

	//开启远程事务
	transationRMSession.Begin()
	//使用mapper
	var activityBean = Activity{
		Id:   "170",
		Name: "rs168-11",
	}
	var updateNum int64 = 0
	var err = exampleActivityMapperImpl.UpdateById(&transationRMSession, activityBean, &updateNum)
	if err != nil {
		panic(err)
	}
	//提交远程事务
	transationRMSession.Commit()
	//关闭远程事务
	//transationRMSession.Rollback()
}

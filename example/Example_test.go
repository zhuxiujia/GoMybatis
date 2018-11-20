package example

import (
	_ "github.com/go-sql-driver/mysql"
	"testing"
	"time"
	"os"
	"fmt"
	"io/ioutil"
	"github.com/zhuxiujia/GoMybatis"
)

//定义mapper文件的接口和结构体，也可以只定义结构体就行
//mapper.go文件 函数必须为2个参数（第一个为自定义结构体参数（属性必须大写），第二个为指针类型的返回数据） error 为返回错误
type ExampleActivityMapperImpl struct {
	SelectAll         func(result *[]Activity) error
	SelectByCondition func(name string, startTime time.Time, endTime time.Time, page int, size int, result *[]Activity) error `mapperParams:"name,startTime,endTime,page,size"`
	UpdateById        func(sessionId string, arg Activity, result *int64) error                                               `mapperParams:"sessionId"` //如果要使用事务，请传入sessionId参数
	Insert            func(arg Activity, result *int64) error
	CountByCondition  func(name string, startTime time.Time, endTime time.Time, result *int) error                            `mapperParams:"name,startTime,endTime"`
}

//初始化mapper文件和结构体
func InitMapper() ExampleActivityMapperImpl {
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
	var exampleActivityMapperImpl = InitMapper()
	//使用mapper
	var result []Activity
	var err = exampleActivityMapperImpl.SelectByCondition("", time.Time{}, time.Time{}, 0, 2000, &result)
	if err != nil {
		panic(err)
	}
	fmt.Println("result=", result)
}

//本地事务使用例子
func Test_local_Transation(t *testing.T) {
	//初始化mapper文件
	exampleActivityMapperImpl := InitMapper()
	//使用事务
	session := *GoMybatis.DefaultSessionFactory.NewSession()
	var sessionId = session.Id()
	session.Begin() //开启事务
	var activityBean = Activity{
		Id:   "170",
		Name: "rs168-6",
	}
	var updateNum int64 = 0
	var e = exampleActivityMapperImpl.UpdateById(sessionId, activityBean, &updateNum)
	fmt.Println("updateNum=", updateNum)
	if e != nil {
		fmt.Println(e)
	}
	session.Commit() //提交事务
	session.Close()  //关闭事务
}

func Test_Remote_Transation(t *testing.T) {

	var addr = "127.0.0.1:17235"
	go GoMybatis.ServerTcp(addr, MysqlDriverName, MysqlUri) //GoMybatis独立事务节点服务器

	time.Sleep(time.Second)

	//初始化mapper文件
	var exampleActivityMapperImpl = InitMapper()

	//使用mapper
	var activityBean = Activity{
		Id:   "170",
		Name: "rs168-4",
	}
	var updateNum int64 = 0
	var err = exampleActivityMapperImpl.UpdateById("", activityBean, &updateNum)
	if err != nil {
		panic(err)
	}
}

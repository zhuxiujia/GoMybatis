package example

import (
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/zhuxiujia/GoMybatis"
	"io/ioutil"
	"reflect"
	"testing"
	"time"
)

//支持基本类型和指针(int,string,time.Time,float...且需要指定参数名称`mapperParams:"name"以逗号隔开，且位置要和实际参数相同)
//参数中包含有*GoMybatis.Session的类型，用于自定义事务
//自定义结构体参数（属性必须大写）
//方法 return 必须包含有error ,为了返回错误信息
type ExampleActivityMapper struct {
	GoMybatis.SessionSupport                                   //session事务操作 写法1.  ExampleActivityMapper.SessionSupport.NewSession()
	NewSession               func() (GoMybatis.Session, error) //session事务操作.写法2   ExampleActivityMapper.NewSession()
	//模板示例
	SelectTemplete      func(name string) ([]Activity, error) `mapperParams:"name"`
	SelectCountTemplete func(name string) (int64, error)      `mapperParams:"name"`
	InsertTemplete      func(arg Activity) (int64, error)
	InsertTempleteBatch func(args []Activity) (int64, error) `mapperParams:"args"`
	UpdateTemplete      func(arg Activity) (int64, error)    `mapperParams:"name"`
	DeleteTemplete      func(name string) (int64, error)     `mapperParams:"name"`

	//传统mybatis示例
	SelectByIds       func(ids []string) ([]Activity, error)       `mapperParams:"ids"`
	SelectByIdMaps    func(ids map[int]string) ([]Activity, error) `mapperParams:"ids"`
	SelectAll         func() ([]map[string]string, error)
	SelectByCondition func(name *string, startTime *time.Time, endTime *time.Time, page *int, size *int) ([]Activity, error) `mapperParams:"name,startTime,endTime,page,size"`
	UpdateById        func(session *GoMybatis.Session, arg Activity) (int64, error)
	Insert            func(arg Activity) (int64, error)
	CountByCondition  func(name string, startTime time.Time, endTime time.Time) (int, error) `mapperParams:"name,startTime,endTime"`
	DeleteById        func(id string) (int64, error)                                         `mapperParams:"id"`
	Choose            func(deleteFlag int) ([]Activity, error)                               `mapperParams:"deleteFlag"`
	SelectLinks       func(column string) ([]Activity, error)                                `mapperParams:"column"`
}

var engine GoMybatis.GoMybatisEngine

//初始化mapper文件和结构体
var exampleActivityMapper = ExampleActivityMapper{}

type TestService struct {
	exampleActivityMapper *ExampleActivityMapper
	UpdateName            func(id string, name string) error `tx:"REQUIRED",rollback:"error"` //事务注解,`tx:"" 开启事务，`tx:"REQUIRED,error"` 指定传播行为为REQUIRED，回滚操作为error类型
	UpdateRemark          func(id string, remark string) error
}

func init() {
	if MysqlUri == "*" {
		println("GoMybatisEngine not init! because MysqlUri is * or MysqlUri is ''")
		return
	}
	engine = GoMybatis.GoMybatisEngine{}.New()
	//mysql链接格式为         用户名:密码@(数据库链接地址:端口)/数据库名称   例如root:123456@(***.mysql.rds.aliyuncs.com:3306)/test
	err := engine.Open("mysql", MysqlUri) //此处请按格式填写你的mysql链接，这里用*号代替
	if err != nil {
		panic(err.Error())
	}

	//动态数据源路由(可选)
	/**
	GoMybatis.Open("mysql", MysqlUri)//添加第二个mysql数据库,请把MysqlUri改成你的第二个数据源链接
	var router = GoMybatis.GoMybatisDataSourceRouter{}.New(func(mapperName string) *string {
		//根据包名路由指向数据源
		if strings.Contains(mapperName, "example.") {
			var url = MysqlUri//第二个mysql数据库,请把MysqlUri改成你的第二个数据源链接
			fmt.Println(url)
			return &url
		}
		return nil
	})
	engine.SetDataSourceRouter(&router)
	**/

	//自定义日志实现(可选)
	/**
		engine.SetLogEnable(true)
		engine.SetLog(&GoMybatis.LogStandard{
			PrintlnFunc: func(messages []byte) {
			},
		})
	    **/

	//注册回调(可选)
	/**
		engine.RegisterCallBack(&GoMybatis.CallBack{
			BeforeExec: func(args []reflect.Value, sqlString *string) {
	         //do something
			},
			BeforeQuery: func(args []reflect.Value, sqlString *string) {
				//do something
			},
			AfterExec: func(args []reflect.Value, sqlString string, result *GoMybatis.Result, err *error) {
				//do something
			},
			AfterQuery: func(args []reflect.Value, sqlString string, result *[]map[string][]byte, err *error) {
				//do something
			},
		})
		**/
	//读取mapper xml文件
	bytes, _ := ioutil.ReadFile("Example_ActivityMapper.xml")
	//设置对应的mapper xml文件
	engine.WriteMapperPtr(&exampleActivityMapper, bytes)
}

//插入
func Test_inset(t *testing.T) {
	if MysqlUri == "" || MysqlUri == "*" {
		fmt.Println("no database url define in MysqlConfig.go , you must set the mysql link!")
		return
	}
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
	var activityBean = Activity{
		Id:   "171",
		Name: "rs168",
	}
	var updateNum, e = exampleActivityMapper.UpdateById(nil, activityBean) //sessionId 有值则使用已经创建的session，否则新建一个session
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
	//使用mapper
	var result, err = exampleActivityMapper.DeleteById("171")
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
	//使用mapper
	name := "注册"

	var result, err = exampleActivityMapper.SelectByCondition(&name, nil, nil, nil, nil)
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
	//使用mapper
	var result, err = exampleActivityMapper.SelectAll()
	if err != nil {
		panic(err)
	}
	var b, _ = json.Marshal(result)
	fmt.Println("result=", string(b))
}

//查询
func Test_count(t *testing.T) {
	if MysqlUri == "" || MysqlUri == "*" {
		fmt.Println("no database url define in MysqlConfig.go , you must set the mysql link!")
		return
	}
	//使用mapper
	var result, err = exampleActivityMapper.CountByCondition("", time.Time{}, time.Time{})
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
	//使用mapper
	var ids = []string{"1", "2"}
	var result, err = exampleActivityMapper.SelectByIds(ids)
	if err != nil {
		panic(err)
	}
	fmt.Println("result=", result)
}

//本地GoMybatis使用例子
func Test_ForEach_Map(t *testing.T) {
	if MysqlUri == "" || MysqlUri == "*" {
		fmt.Println("no database url define in MysqlConfig.go , you must set the mysql link!")
		return
	}
	//使用mapper
	var ids = map[int]string{1: "165", 2: "166"}
	var result, err = exampleActivityMapper.SelectByIdMaps(ids)
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
	//使用事务
	var session, err = exampleActivityMapper.SessionSupport.NewSession()
	if err != nil {
		t.Fatal(err)
	}
	session.Begin() //开启事务
	var activityBean = Activity{
		Id:   "170",
		Name: "rs168-8",
	}
	var updateNum, e = exampleActivityMapper.UpdateById(&session, activityBean) //sessionId 有值则使用已经创建的session，否则新建一个session
	fmt.Println("updateNum=", updateNum)
	if e != nil {
		panic(e)
	}
	session.Commit() //提交事务
	session.Close()  //关闭事务
}

func Test_choose(t *testing.T) {
	if MysqlUri == "" || MysqlUri == "*" {
		fmt.Println("no database url define in MysqlConfig.go , you must set the mysql link!")
		return
	}
	//使用mapper
	var result, err = exampleActivityMapper.Choose(1)
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
	//使用mapper
	var result, err = exampleActivityMapper.SelectLinks("name")
	if err != nil {
		panic(err)
	}
	fmt.Println("result=", result)
}

func TestSelectTemplete(t *testing.T) {
	if MysqlUri == "" || MysqlUri == "*" {
		fmt.Println("no database url define in MysqlConfig.go , you must set the mysql link!")
		return
	}
	//使用mapper
	var result, err = exampleActivityMapper.SelectTemplete("hello")
	if err != nil {
		panic(err)
	}
	fmt.Println("result=", result)
}

func TestSelectCountTemplete(t *testing.T) {
	if MysqlUri == "" || MysqlUri == "*" {
		fmt.Println("no database url define in MysqlConfig.go , you must set the mysql link!")
		return
	}
	//使用mapper
	var result, err = exampleActivityMapper.SelectCountTemplete("hello")
	if err != nil {
		panic(err)
	}
	fmt.Println("result=", result)
}

func TestInsertTemplete(t *testing.T) {
	if MysqlUri == "" || MysqlUri == "*" {
		fmt.Println("no database url define in MysqlConfig.go , you must set the mysql link!")
		return
	}
	//使用mapper
	var result, err = exampleActivityMapper.InsertTemplete(Activity{Id: "178", Name: "test_insret", CreateTime: time.Now(), DeleteFlag: 1})
	if err != nil {
		panic(err)
	}
	fmt.Println("result=", result)
}

//批量插入模板
func TestInsertTempleteBatch(t *testing.T) {
	if MysqlUri == "" || MysqlUri == "*" {
		fmt.Println("no database url define in MysqlConfig.go , you must set the mysql link!")
		return
	}
	var args = []Activity{
		{
			Id:   "221",
			Name: "test",
		},
		{
			Id:   "222",
			Name: "test",
		},
		{
			Id:   "223",
			Name: "test",
		},
	}
	n, err := exampleActivityMapper.InsertTempleteBatch(args)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("updateNum", n)
	time.Sleep(time.Second)
}

//修改模板默认支持逻辑删除和乐观锁
func TestUpdateTemplete(t *testing.T) {
	if MysqlUri == "" || MysqlUri == "*" {
		fmt.Println("no database url define in MysqlConfig.go , you must set the mysql link!")
		return
	}
	var activityBean = Activity{
		Id:      "171",
		Name:    "rs168",
		Version: 2,
	}
	//会自动生成乐观锁和逻辑删除字段 set version= * where version = * and delete_flag = *
	// update set name = 'rs168',version = 1 from biz_activity where name = 'rs168' and delete_flag = 1 and version = 0
	var updateNum, e = exampleActivityMapper.UpdateTemplete(activityBean)
	fmt.Println("updateNum=", updateNum)
	if e != nil {
		panic(e)
	}
}

//删除
func TestDeleteTemplete(t *testing.T) {
	if MysqlUri == "" || MysqlUri == "*" {
		fmt.Println("no database url define in MysqlConfig.go , you must set the mysql link!")
		return
	}
	//模板默认支持逻辑删除
	var result, err = exampleActivityMapper.DeleteTemplete("rs168")
	if err != nil {
		panic(err)
	}
	fmt.Println("result=", result)
}

//嵌套事务/带有传播行为的事务
func TestTestService(t *testing.T) {
	if MysqlUri == "" || MysqlUri == "*" {
		fmt.Println("no database url define in MysqlConfig.go , you must set the mysql link!")
		return
	}
	var testService = initTestService()

	//go testService.UpdateName("167", "updated name1")
	testService.UpdateName("167", "updated name2")

	time.Sleep(5 * time.Second)
}

func initTestService() TestService {
	var testService  TestService
	testService = TestService{
		exampleActivityMapper: &exampleActivityMapper,
		UpdateRemark: func(id string, remark string) error {
			println("UpdateRemark start")
			return nil
		},
		UpdateName: func(id string, name string) error {
			println("UpdateName start")
			var activitys, err = testService.exampleActivityMapper.SelectByIds([]string{id})
			if err != nil {
				panic(err)
			}
			var activity = activitys[0]
			activity.Name = name
			updateNum, err := testService.exampleActivityMapper.UpdateTemplete(activity)
			if err != nil {
				panic(err)
			}
			println("success updateNum:", updateNum)
			testService.UpdateRemark(id, "updated remark")
			return nil
		},
	}
	GoMybatis.AopProxyService(reflect.ValueOf(&testService), &engine)
	return testService
}

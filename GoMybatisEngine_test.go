package GoMybatis

import (
	"fmt"
	"github.com/zhuxiujia/GoMybatis/example"
	"github.com/zhuxiujia/GoMybatis/utils"
	"sync"
	"testing"
	"time"
)

//假设Mysql 数据库查询时间为0，框架单协程的Benchmark性能
func Benchmark_One_Transcation(b *testing.B) {
	b.StopTimer()
	//使用事务
	session := Session(&TestSession{})
	//初始化mapper文件
	var exampleActivityMapperImpl = InitMapperByLocalSession()
	//使用mapper
	//开始压力测试
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		var _, err = exampleActivityMapperImpl.SelectByCondition(&session, nil, nil, nil, nil, nil)
		if err != nil {
			b.Fatal(err)
		}
		//fmt.Println(data)
	}
}

//假设Mysql 数据库查询时间为0，框架多个协程的并发数的性能
func Benchmark_One_Transcation_multiple_coroutine(b *testing.B) {
	b.StopTimer()
	//使用事务
	session := Session(&TestSession{})
	//初始化mapper文件
	var exampleActivityMapperImpl = InitMapperByLocalSession()
	//使用mapper
	////开始TPS测试
	var total = b.N      //总数
	var goruntine = 1000 //并发数
	if total <= 1 {
		goruntine = 1
	}
	var waitGroup = sync.WaitGroup{}
	waitGroup.Add(goruntine)

	b.StartTimer()
	for i := 0; i < goruntine; i++ {
		go func() {
			var itemCount = total / goruntine
			for f := 0; f < itemCount; f++ {
				_, e := exampleActivityMapperImpl.SelectByCondition(&session, nil, nil, nil, nil, nil)
				if e != nil {
					b.Fatal(e)
				}
			}
			waitGroup.Done()
		}()
	}
	waitGroup.Wait()
}

//假设Mysql 数据库查询时间为0，框架单协程的并发数的性能
func Test_One_Transcation_TPS(t *testing.T) {
	//使用事务
	session := Session(&TestSession{})
	//初始化mapper文件
	var exampleActivityMapperImpl = InitMapperByLocalSession()
	//使用mapper
	var name = "dsa"
	var times = time.Now()
	var page = 1
	//开始TPS测试
	var total = 10000
	defer utils.CountMethodTps(float64(total), time.Now(), "Test_One_Transcation_TPS")
	for i := 0; i < total; i++ {
		var _, err = exampleActivityMapperImpl.SelectByCondition(&session, &name, &times, &times, &page, &page)
		if err != nil {
			t.Fatal(err)
		}
	}
}

//假设Mysql 数据库查询时间为0，框架多个协程的并发数的性能
func Test_One_Transcation_multiple_coroutine_TPS(t *testing.T) {
	//使用事务
	session := Session(&TestSession{})
	//初始化mapper文件
	var exampleActivityMapperImpl = InitMapperByLocalSession()
	//使用mapper

	var name = "dsa"
	var times = time.Now()
	var page = 1

	////开始TPS测试
	var total = 100000   //总数
	var goruntine = 1000 //并发数
	var waitGroup = sync.WaitGroup{}
	waitGroup.Add(goruntine)

	defer utils.CountMethodTps(float64(total), time.Now(), "Test_One_Transcation_multiple_coroutine_TPS")

	for i := 0; i < goruntine; i++ {
		go func() {
			var itemCount = total / goruntine
			for f := 0; f < itemCount; f++ {
				_, e := exampleActivityMapperImpl.SelectByCondition(&session, &name, &times, &times, &page, &page)
				if e != nil {
					t.Fatal(e)
				}
			}
			waitGroup.Done()
		}()
	}
	waitGroup.Wait()
}

//验证测试直接返回数据
func Test_Transcation(t *testing.T) {
	//使用事务
	session := Session(&TestSession{})
	//初始化mapper文件
	var exampleActivityMapperImpl = InitMapperByLocalSession()
	//使用mapper
	var results, err = exampleActivityMapperImpl.SelectByCondition(&session, nil, nil, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(results)
}

type TestSession struct {
	Session
}

func (it *TestSession) Id() string {
	return "sadf"
}
func (it *TestSession) Query(sqlorArgs string) ([]map[string][]byte, error) {
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
func (it *TestSession) Exec(sqlorArgs string) (*Result, error) {
	return nil, nil
}
func (it *TestSession) Rollback() error {
	return nil
}
func (it *TestSession) Commit() error {
	return nil
}
func (it *TestSession) Begin() error {
	return nil
}
func (it *TestSession) Close() {

}

//定义mapper文件的接口和结构体
// 支持基本类型(int,string,time.Time,float...且需要指定参数名称`mapperParams:"name"以逗号隔开，且位置要和实际参数相同)
//自定义结构体参数（属性必须大写）
//参数中除了session指针外，为指针类型的皆为数据
// 函数return必须为error 为返回错误信息
type ExampleActivityMapperImpl struct {
	SelectByCondition func(session *Session, name *string, startTime *time.Time, endTime *time.Time, page *int, size *int) ([]example.Activity, error) `mapperParams:"session,name,startTime,endTime,page,size"`
}

//初始化mapper文件和结构体
func InitMapperByLocalSession() ExampleActivityMapperImpl {
	var engine = GoMybatisEngine{}.New()
	//mysql链接格式为         用户名:密码@(数据库链接地址:端口)/数据库名称   例如root:123456@(***.mysql.rds.aliyuncs.com:3306)/test
	err := engine.Open("mysql", "") //此处请按格式填写你的mysql链接，这里用*号代替
	if err != nil {
		panic(err.Error())
	}

	var bytes = []byte(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE mapper PUBLIC "-//mybatis.org//DTD Mapper 3.0//EN"
        "https://raw.githubusercontent.com/zhuxiujia/GoMybatis/master/mybatis-3-mapper.dtd">
<mapper>
    <resultMap id="BaseResultMap" type="example.Activity">
        <id column="id" property="id"/>
        <result column="name" property="name" langType="string"/>
        <result column="pc_link" property="pcLink" langType="string"/>
        <result column="h5_link" property="h5Link" langType="string"/>
        <result column="remark" property="remark" langType="string"/>
        <result column="create_time" property="createTime" langType="time.Time"/>
        <result column="delete_flag" property="deleteFlag" langType="int"/>
    </resultMap>

    <!--List<Activity> selectByCondition(@Param("name") String name,@Param("startTime") Date startTime,@Param("endTime") Date endTime,@Param("index") Integer index,@Param("size") Integer size);-->
    <!-- 后台查询产品 -->
    <select id="selectByCondition" resultMap="BaseResultMap">
        select * from biz_activity where delete_flag=1
        <if test="name != nil">
            and name like #{name}
        </if>
        <if test="startTime != nil">
            and create_time >= #{startTime}
        </if>
        <if test="endTime != nil">
            and create_time &lt;= #{endTime}
        </if>
        order by create_time desc
        <if test="page != nil and size != nil">limit #{page}, #{size}</if>
    </select>`)
	var exampleActivityMapperImpl ExampleActivityMapperImpl
	//设置对应的mapper xml文件,禁止输出日志
	engine.SetLogEnable(false)
	engine.WriteMapperPtr(&exampleActivityMapperImpl, bytes)
	return exampleActivityMapperImpl
}

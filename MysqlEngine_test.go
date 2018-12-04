package GoMybatis

import (
	"time"
	"github.com/zhuxiujia/GoMybatis/example"
	"testing"
	"github.com/zhuxiujia/GoMybatis/utils"
)

//假设Mysql 数据库查询时间为0，框架的并发数的性能
func Test_One_Transcation_TPS(t *testing.T) {
	//使用事务
	session := Session(&TestSession{})
	//初始化mapper文件
	var exampleActivityMapperImpl = InitMapperByLocalSession()
	//使用mapper

	//开始TPS测试
	var total = 100000
	defer utils.CountMethodTps(float64(total), time.Now(), "Test_One_Transcation_TPS")
	var results []example.Activity
	for i := 0; i < total; i++ {
		var err = exampleActivityMapperImpl.SelectByCondition(&session, "", time.Time{}, time.Time{}, 0, 2000, &results)
		if err != nil {
			panic(err)
		}
	}
}

//验证测试直接返回数据
func Test_Transcation(t *testing.T) {
	//使用事务
	session := Session(&TestSession{})
	//初始化mapper文件
	var exampleActivityMapperImpl = InitMapperByLocalSession()
	//使用mapper

	//开始TPS测试
	var results []example.Activity
	var err = exampleActivityMapperImpl.SelectByCondition(&session, "", time.Time{}, time.Time{}, 0, 2000, &results)
	if err != nil {
		panic(err)
	}
}

type TestSession struct {
	Session
}

func (this *TestSession) Id() string {
	return "sadf"
}
func (this *TestSession) Query(sqlorArgs string) ([]map[string][]byte, error) {
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
func (this *TestSession) Exec(sqlorArgs string) (*Result, error) {
	return nil, nil
}
func (this *TestSession) Rollback() error {
	return nil
}
func (this *TestSession) Commit() error {
	return nil
}
func (this *TestSession) Begin() error {
	return nil
}
func (this *TestSession) Close() {

}

//定义mapper文件的接口和结构体
// 支持基本类型(int,string,time.Time,float...且需要指定参数名称`mapperParams:"name"以逗号隔开，且位置要和实际参数相同)
//自定义结构体参数（属性必须大写）
//参数中除了session指针外，为指针类型的皆为返回数据
// 函数return必须为error 为返回错误信息
type ExampleActivityMapperImpl struct {
	SelectByCondition func(session *Session, name string, startTime time.Time, endTime time.Time, page int, size int, result *[]example.Activity) error `mapperParams:"session,name,startTime,endTime,page,size"`
}

//初始化mapper文件和结构体
func InitMapperByLocalSession() ExampleActivityMapperImpl {
	var err error
	//mysql链接格式为         用户名:密码@(数据库链接地址:端口)/数据库名称   例如root:123456@(***.mysql.rds.aliyuncs.com:3306)/test
	engine, err := Open("mysql", "") //此处请按格式填写你的mysql链接，这里用*号代替
	if err != nil {
		panic(err.Error())
	}

	var bytes = []byte(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE mapper PUBLIC "-//mybatis.org//DTD Mapper 3.0//EN"
        "https://raw.githubusercontent.com/zhuxiujia/GoMybatis/master/mybatis-3-mapper.dtd">
<mapper>
    <resultMap id="BaseResultMap" type="example.Activity">
        <id column="id" property="id"/>
        <result column="name" property="name" goType="string"/>
        <result column="pc_link" property="pcLink" goType="string"/>
        <result column="h5_link" property="h5Link" goType="string"/>
        <result column="remark" property="remark" goType="string"/>
        <result column="create_time" property="createTime" goType="time.Time"/>
        <result column="delete_flag" property="deleteFlag" goType="int"/>
    </resultMap>

    <!--List<Activity> selectByCondition(@Param("name") String name,@Param("startTime") Date startTime,@Param("endTime") Date endTime,@Param("index") Integer index,@Param("size") Integer size);-->
    <!-- 后台查询产品 -->
    <select id="selectByCondition" resultMap="BaseResultMap">
        select * from biz_activity where delete_flag=1
        <bind name="pattern" value="'%' + name + '%'" />
        <if test="name != ''">
            <!--可以使用bind标签 and name like #{pattern}-->
            and name like #{pattern}
            <!--可以使用默认 and name like concat('%',#{name},'%')-->
            <!--and name like concat('%',#{name},'%')-->
        </if>
        <if test="startTime != ''">
            and create_time >= #{startTime}
        </if>
        <if test="endTime != ''">
            and create_time &lt;= #{endTime}
        </if>
        order by create_time desc
        <if test="page >= 0 and size != 0">limit #{page}, #{size}</if>
    </select>`)
	var exampleActivityMapperImpl ExampleActivityMapperImpl
	//设置对应的mapper xml文件,禁止输出日志
	UseProxyMapperByEngine(&exampleActivityMapperImpl, bytes, engine,false)
	return exampleActivityMapperImpl
}
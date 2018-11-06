# GoMybatis
# 文档网站站点
https://zhuxiujia.github.io/gomybatis.io/#/
GoMybatis 是根据java版 Mybatis3 的实现,基于Xorm的Engine和govaluate表达式及反射实现。
GoMybatis 内部在初始化时反射分析mapper xml生成golang的func代码，默认支持绝大部分的Java版的mybatis标签和规范,
支持标签
`<select>,<update>,<insert>,<delete>,<trim>,<if>,<set>,<foreach>`
# 使用教程
<pre>
go get github.com/zhuxiujia/GoMybatis
</pre>
mapper.go 文件案例
<pre>
//属性必须大写,GoBatis将会使用反射取得字段名称和值
type SelectByConditionArg struct {
	Name      string
	StartTime time.Time
	EndTime   time.Time
	Page      int
	Size      int
}
type ActivityMapperImpl struct {
  //mapper.go文件 函数必须为2个参数（第一个为自定义结构体参数（属性必须大写），第二个为指针类型的返回数据） error 为返回错误
	SelectByCondition func(arg SelectByConditionArg, result *[]model.Activity) error
}
</pre>

xml文件案例:
```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE mapper PUBLIC "-//mybatis.org//DTD Mapper 3.0//EN" "http://mybatis.org/dtd/mybatis-3-mapper.dtd">
<mapper namespace="ActivityMapperImpl">
    <resultMap id="BaseResultMap" type="model.Activity">
        <id column="id" property="id" jdbcType="VARCHAR"/>
        <result column="name" property="name" jdbcType="VARCHAR"/>
        <result column="pc_link" property="pcLink" jdbcType="VARCHAR"/>
        <result column="h5_link" property="h5Link" jdbcType="VARCHAR"/>
        <result column="remark" property="remark" jdbcType="VARCHAR"/>
        <result column="create_time" property="createTime" jdbcType="TIMESTAMP"/>
        <result column="delete_flag" property="deleteFlag" jdbcType="INTEGER"/>
    </resultMap>
    <!-- SelectByCondition func(arg SelectByConditionArg, result *[]model.Activity) error -->
    <!-- 后台查询产品 -->
    <select id="SelectByCondition" resultMap="BaseResultMap">
        select
        <trim prefix="" suffix="" suffixOverrides=",">
            <if test="Name != ''">name,</if>
        </trim>
        from biz_activity where delete_flag=1
        <if test="Name != ''">
            and name like concat('%',#{Name},'%')
        </if>
        <if test="StartTime != 0">
            and create_time >= #{StartTime}
        </if>
        <if test="EndTime != 0">
            and create_time &lt;= #{EndTime}
        </if>
        order by create_time desc
        <if test="Page != 0 and Size != 0">limit #{Page}, #{Size}</if>
    </select>
</mapper>
```
在服务层实际使用mapper
<pre>
import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/Knetic/govaluate"
	)
	
func main() {
  var mapper ActivityMapperImpl
  engine, dbError := xorm.NewEngine("mysql", "")
	if dbError != nil {
		fmt.Println(dbError)
		return
	}
  engine.LogMode(true)
  UseProxyMapper(&mapper, engine.DB)
  //查询
  var r []model.Activity //model.Activity 此处应改为你自己的数据库模型类型
  var err = mapper.SelectByCondition(SelectByConditionArg{
		Name: `rs`,
	}, &r)
	fmt.Println(err)
	fmt.Println(r)
}
</pre>


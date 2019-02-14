package GoMybatis

import (
	"fmt"
	"github.com/zhuxiujia/GoMybatis/example"
	"github.com/zhuxiujia/GoMybatis/lib/github.com/Knetic/govaluate"
	"github.com/zhuxiujia/GoMybatis/utils"
	"testing"
	"time"
)

//压力测试 sql构建情况
func Benchmark_SqlBuilder(b *testing.B) {
	b.StopTimer()
	var mapper = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE mapper PUBLIC "-//mybatis.org//DTD Mapper 3.0//EN" "http://mybatis.org/dtd/mybatis-3-mapper.dtd">
<mapper>
    <!--List<Activity> selectByCondition(@Param("name") String name,@Param("startTime") Date startTime,@Param("endTime") Date endTime,@Param("index") Integer index,@Param("size") Integer size);-->
    <!-- 后台查询产品 -->
    <select id="selectByCondition">
        select * from biz_activity where delete_flag=1
        <if test="name != ''">
            and name like concat('%',#{name},'%')
        </if>
        <if test="startTime != ''">
            and create_time >= #{startTime}
        </if>
        <if test="endTime != ''">
            and create_time &lt;= #{endTime}
        </if>
        order by create_time desc
        <if test="page != '' and size != ''">limit #{page}, #{size}</if>
    </select>
</mapper>`
	var mapperTree = LoadMapperXml([]byte(mapper))

	var builder = GoMybatisSqlBuilder{}.New(GoMybatisExpressionTypeConvert{}, GoMybatisSqlArgTypeConvert{}, ExpressionEngineProxy{}.New(&ExpressionEngineExpr{}, true), &LogStandard{}, false)
	var paramMap = make(map[string]interface{})
	paramMap["name"] = "sda"
	paramMap["startTime"] = ""
	paramMap["endTime"] = ""
	paramMap["page"] = "1"
	paramMap["size"] = "20"

	//paramMap["type_name"] =  StringType
	//paramMap["type_startTime"] = StringType
	//paramMap["type_endTime"] = StringType
	//paramMap["type_page"] = IntType
	//paramMap["type_size"] = IntType

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		builder.BuildSql(paramMap, mapperTree["selectByCondition"])
	}
}

//测试sql生成tps
func Test_SqlBuilder_Tps(t *testing.T) {
	var mapper = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE mapper PUBLIC "-//mybatis.org//DTD Mapper 3.0//EN" "http://mybatis.org/dtd/mybatis-3-mapper.dtd">
<mapper>
    <!--List<Activity> selectByCondition(@Param("name") String name,@Param("startTime") Date startTime,@Param("endTime") Date endTime,@Param("index") Integer index,@Param("size") Integer size);-->
    <!-- 后台查询产品 -->
    <select id="selectByCondition">
        select * from biz_activity where delete_flag=1
        <if test="name != nil">
            and name like concat('%',#{name},'%')
        </if>
        <if test="startTime != nil">
            and create_time >= #{startTime}
        </if>
        <if test="endTime != nil">
            and create_time &lt;= #{endTime}
        </if>
        order by create_time desc
        <if test="page >= 0 and size != 0">limit #{page}, #{size}</if>
    </select>
</mapper>`
	var mapperTree = LoadMapperXml([]byte(mapper))

	var builder = GoMybatisSqlBuilder{}.New(GoMybatisExpressionTypeConvert{}, GoMybatisSqlArgTypeConvert{}, ExpressionEngineProxy{}.New(&ExpressionEngineExpr{}, true), &LogStandard{}, false)
	var paramMap = make(map[string]interface{})
	paramMap["name"] = ""
	paramMap["startTime"] = ""
	paramMap["endTime"] = ""
	paramMap["page"] = 0
	paramMap["size"] = 0

	var startTime = time.Now()
	for i := 0; i < 100000; i++ {
		//var sql, e =
		builder.BuildSql(paramMap, mapperTree["selectByCondition"])
		//fmt.Println(sql, e)
	}
	utils.CountMethodTps(100000, startTime, "Test_SqlBuilder_Tps")
}

func Test_reflect_tps(t *testing.T) {
	var p = make(map[string]string)
	var n = p
	n["a"] = "b"
	fmt.Println(p)

	defer utils.CountMethodTps(10000, time.Now(), "Test_reflect_tps")

	for k := 0; k < 10000; k++ {
		evalExpression, _ := govaluate.NewEvaluableExpression("name != nil")
		//fmt.Println(err)
		var p = make(map[string]interface{})
		p["name"] = "sdaf"
		evalExpression.Evaluate(p)
		//fmt.Println(err)
		//fmt.Println(result)
	}

}

func Test_bind_string(t *testing.T) {
	var activity = example.Activity{
		Id:         "1",
		DeleteFlag: 1,
	}
	var evaluateParameters = make(map[string]interface{})
	evaluateParameters["activity"] = activity
	var expression = "'%' + activity.Id + '%'"
	evalExpression, err := govaluate.NewEvaluableExpression(expression)
	if err != nil {
		t.Fatal(err)
	}
	result, err := evalExpression.Evaluate(evaluateParameters)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(result)
}

func TestGoMybatisSqlBuilder_BuildSql(t *testing.T) {
	var mapper = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE mapper PUBLIC "-//mybatis.org//DTD Mapper 3.0//EN" "http://mybatis.org/dtd/mybatis-3-mapper.dtd">
<mapper>
    <resultMap id="BaseResultMap">
        <id column="id" property="id"/>
        <result column="name" property="name" langType="string"/>
        <result column="pc_link" property="pcLink" langType="string"/>
        <result column="h5_link" property="h5Link" langType="string"/>
        <result column="remark" property="remark" langType="string"/>
        <result column="create_time" property="createTime" langType="time.Time"/>
        <result column="delete_flag" property="deleteFlag" langType="int"/>
    </resultMap>
    <select id="selectByCondition" resultMap="BaseResultMap">
        <bind name="pattern" value="'%' + name + '%'"/>
        select * from biz_activity
        <where>
            <if test="name != nil">
                and name like #{pattern}
            </if>
            <if test="startTime != nil">and create_time >= #{startTime}</if>
            <if test="endTime != nil">and create_time &lt;= #{endTime}</if>
        </where>
        order by 
        <trim prefix="" suffix="" suffixOverrides=",">
            <if test="name != nil">name,</if>
        </trim>
        desc
        <choose>
            <when test="page < 1">limit 3</when>
            <when test="page > 1">limit 2</when>
            <otherwise>limit 1</otherwise>
        </choose>
    </select>
</mapper>`
	var mapperTree = LoadMapperXml([]byte(mapper))

	var builder = GoMybatisSqlBuilder{}.New(GoMybatisExpressionTypeConvert{}, GoMybatisSqlArgTypeConvert{}, ExpressionEngineProxy{}.New(&ExpressionEngineExpr{}, true), &LogStandard{}, true)
	var paramMap = make(map[string]interface{})
	paramMap["name"] = "name"
	paramMap["type_name"] = StringType
	paramMap["startTime"] = nil
	paramMap["endTime"] = nil
	paramMap["page"] = 0
	paramMap["size"] = 0
	var sql, err = builder.BuildSql(paramMap, mapperTree["selectByCondition"])
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(sql)
}

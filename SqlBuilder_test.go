package GoMybatis

import (
	"fmt"
	"github.com/zhuxiujia/GoMybatis/lib/github.com/Knetic/govaluate"
	"github.com/zhuxiujia/GoMybatis/utils"
	"reflect"
	"testing"
	"time"
)

//测试sql生成tps
func Test_SqlBuilder_Tps2(t *testing.T) {
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
        <if test="page >= 0 and size != 0">limit #{page}, #{size}</if>
    </select>
</mapper>`
	var mapperTree = LoadMapperXml([]byte(mapper))

	var builder = GoMybatisSqlBuilder{}.New(GoMybatisExpressionTypeConvert{}, GoMybatisSqlArgTypeConvert{})
	var paramMap = make(map[string]SqlArg)
	paramMap["name"] = SqlArg{
		Value: "",
		Type:  reflect.TypeOf(""),
	}
	paramMap["startTime"] = SqlArg{
		Value: "",
		Type:  reflect.TypeOf(""),
	}
	paramMap["endTime"] = SqlArg{
		Value: "",
		Type:  reflect.TypeOf(""),
	}
	paramMap["page"] = SqlArg{
		Value: 0,
		Type:  reflect.TypeOf(0),
	}
	paramMap["size"] = SqlArg{
		Value: 0,
		Type:  reflect.TypeOf(0),
	}
	defer utils.CountMethodTps(100000, time.Now(), "Test_SqlBuilder_Tps")
	for i := 0; i < 100000; i++ {
		//var sql, e =
		builder.BuildSql(paramMap, mapperTree["selectByCondition"])
		//fmt.Println(sql, e)
	}
}

func Test_reflect_tps(t *testing.T) {
	var p = make(map[string]string)
	var n = p
	n["a"] = "b"
	fmt.Println(p)

	defer utils.CountMethodTps(100000, time.Now(), "Test_reflect_tps")

	for k := 0; k < 100000; k++ {
		evalExpression, _ := govaluate.NewEvaluableExpression("name != ''")
		//fmt.Println(err)
		var p = make(map[string]interface{})
		p["name"] = "sdaf"
		evalExpression.Evaluate(p)
		//fmt.Println(err)
		//fmt.Println(result)
	}

}

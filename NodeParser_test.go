package GoMybatis

import (
	"fmt"
	"testing"
)

func TestNodeParser_ParserNodes(t *testing.T) {
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
	//fmt.Println(mapperTree)

	var nodePar = NodeParser{}
	var sqlNodes = nodePar.ParserNodes(mapperTree["selectByCondition"].ElementItems)

	fmt.Println(sqlNodes)

	var proxy = ExpressionEngineProxy{}.New(&ExpressionEngineGoExpress{}, true)

	var argMap = map[string]interface{}{
		"SqlArgTypeConvert":      &GoMybatisSqlArgTypeConvert{},
		"*ExpressionEngineProxy": &proxy,
		"name":                   "sadf",
	}
	argMap["name"] = ""
	argMap["startTime"] = ""
	argMap["endTime"] = ""
	argMap["page"] = 0
	argMap["size"] = 0
	var r, e = DoChildNodes(sqlNodes, argMap)
	if e != nil {
		t.Fatal(e)
	}
	fmt.Println(string(r))
}

func BenchmarkNodeParser_ParserNodes(b *testing.B) {
	b.StopTimer()
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
	//fmt.Println(mapperTree)

	var nodePar = NodeParser{}
	var sqlNodes = nodePar.ParserNodes(mapperTree["selectByCondition"].ElementItems)
	var proxy = ExpressionEngineProxy{}.New(&ExpressionEngineGoExpress{}, true)

	var argMap = map[string]interface{}{
		"SqlArgTypeConvert":      &GoMybatisSqlArgTypeConvert{},
		"*ExpressionEngineProxy": &proxy,
		"name":                   "sadf",
	}
	argMap["name"] = ""
	argMap["startTime"] = ""
	argMap["endTime"] = ""
	argMap["page"] = 0
	argMap["size"] = 0

	b.StartTimer()
	for i := 0; i < b.N; i++ {

		if true {
			var b = []byte(argMap["name"].(string))
			if b != nil {

			}
			continue
		}
		var _, e = DoChildNodes(sqlNodes, argMap)
		if e != nil {
			b.Fatal(e)
		}
		//fmt.Println(r)
	}

}

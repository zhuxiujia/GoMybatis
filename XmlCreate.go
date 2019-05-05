package GoMybatis

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

var XmlData = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE mapper PUBLIC "-//mybatis.org//DTD Mapper 3.0//EN"
        "https://raw.githubusercontent.com/zhuxiujia/GoMybatis/master/mybatis-3-mapper.dtd">
<mapper>
    <!--logic_enable 逻辑删除字段-->
    <!--logic_deleted 逻辑删除已删除字段-->
    <!--logic_undelete 逻辑删除 未删除字段-->
    <!--version_enable 乐观锁版本字段,支持int,int8,int16,int32,int64-->
    <resultMap id="BaseResultMap" tables="#{table}">
    #{resultMapBody}
    </resultMap>
</mapper>
`
var XmlLogicEnable = `logic_enable="true" logic_undelete="1" logic_deleted="0"`
var XmlVersionEnable = `version_enable="true"`
var XmlIdItem = `<id column="id" property="id"/>`
var ResultItem = `<result column="#{property}" property="#{property}" langType="#{langType}" #{version} #{logic} />`

/**
//例子

func TestUserAddres(t *testing.T)  {
	var s=utils.CreateDefaultXml("biz_user_address",example.Activity{})//创建xml内容
	utils.WriteXml("D:/GOPATH/src/dao/ActivityMapper.xml",[]byte(s))//写入磁盘
}
 */
//根据结构体 创建xml文件
func CreateXml(tableName string, bean interface{}) []byte {
	var content = ""
	var tv = reflect.TypeOf(bean)
	if tv.Kind() == reflect.Ptr {
		tv = tv.Elem()
	}
	for i := 0; i < tv.NumField(); i++ {
		var item = tv.Field(i)
		if item.Name == "id" {
			content += XmlIdItem
			content += "\n"
		} else {
			var itemName = item.Name
			var itemJson = item.Tag.Get("json")
			if itemJson != "" {
				itemName = itemJson
			}

			var itemStr = strings.Replace(ResultItem, "#{property}", itemName, -1)
			itemStr = strings.Replace(itemStr, "#{langType}", item.Type.Name(), -1)

			var gm = item.Tag.Get("gm")
			if gm != "" {
				if gm == "id" {
					content += XmlIdItem
					content += "\n"
				}
				if gm == "version" {
					itemStr = strings.Replace(itemStr, "#{version}", XmlVersionEnable, -1)
				}
				if gm == "logic" {
					itemStr = strings.Replace(itemStr, "#{logic}", XmlLogicEnable, -1)
				}
			}
			//clean
			itemStr = strings.Replace(itemStr, "#{version}", "", -1)
			itemStr = strings.Replace(itemStr, "#{logic}", "", -1)

			content += "\t" + itemStr
			if i+1 < tv.NumField() {
				content += "\n"
			}
		}
	}
	var res = strings.Replace(XmlData, "#{resultMapBody}", content, -1)
	res = strings.Replace(res, "#{table}", tableName, -1)
	return []byte(res)
}

//写文件到当前路径
func WriteXml(fileName string, body []byte) {
	f, err := os.Create(fileName)
	defer f.Close()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		_, err = f.Write(body)
		if err != nil {
			println(err)
		} else {
			println("写入文件成功：" + fileName)
		}
	}
}

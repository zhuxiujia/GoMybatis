# SQL mapper framework for Golang
[![Go Report Card](https://goreportcard.com/badge/github.com/zhuxiujia/GoMybatis)](https://goreportcard.com/report/github.com/zhuxiujia/GoMybatis)
[![Build Status](https://travis-ci.com/zhuxiujia/GoMybatis.svg?branch=master)](https://travis-ci.com/zhuxiujia/GoMybatis)
[![GoDoc](https://godoc.org/github.com/zhuxiujia/GoMybatis?status.svg)](https://godoc.org/github.com/zhuxiujia/GoMybatis)
[![Coverage Status](https://coveralls.io/repos/github/zhuxiujia/GoMybatis/badge.svg?branch=master)](https://coveralls.io/github/zhuxiujia/GoMybatis?branch=master)
[![codecov](https://codecov.io/gh/zhuxiujia/GoMybatis/branch/master/graph/badge.svg)](https://codecov.io/gh/zhuxiujia/GoMybatis)


![Image text](https://zhuxiujia.github.io/gomybatis.io/assets/vuetify.png)
# 官方网站/文档
https://zhuxiujia.github.io/gomybatis.io/info.html
# 优势
<a href="https://zhuxiujia.github.io/gomybatis.io/info.html">-协程高并发</a>，假设数据库响应时间为0，在6核16GB PC上框架每秒事务数可达246982Tps/s,耗时仅仅0.4s<br>
<a href="https://zhuxiujia.github.io/gomybatis.io/info.html">-事务支持</a>，session灵活插拔，兼容过渡期微服务<br>
<a href="https://zhuxiujia.github.io/gomybatis.io/info.html">-动态SQL</a>，在xml中可灵活运用if判断，foreach遍历数组，resultMap,bind等等java框架Mybatis包含的实用功能`<select>,<update>,<insert>,<delete>,<trim>,<if>,<set>,<where>,<foreach>,<resultMap>,<bind>,<choose><when><otherwise>,<sql><include>`<br>
<a href="https://zhuxiujia.github.io/gomybatis.io/info.html">-多数据库支持Mysql,Postgres,Tidb,SQLite,Oracle....等等更多</a><br>
<a href="https://zhuxiujia.github.io/gomybatis.io/info.html">-快速上手</a>基于反射动态代理,无需go generate生成*.go等中间代码，xml读取后可直接调用函数<br>
<a href="https://zhuxiujia.github.io/gomybatis.io/info.html">-接口化设计扩展性好</a>面向接口及设计模式，扩展性和替换性好<br>
### 通过远程代理处理事务支持  处于 单数据库(Mysql,postgresql)-分布式数据库（TiDB,cockroachdb...）过渡期间的微服务
![Image text](https://zhuxiujia.github.io/gomybatis.io/assets/tx.png)


数据库驱动列表
```
 Mysql: github.com/go-sql-driver/mysql
 MyMysql: github.com/ziutek/mymysql/godrv
 Postgres: github.com/lib/pq
 Tidb: github.com/pingcap/tidb
 SQLite: github.com/mattn/go-sqlite3
 MsSql: github.com/denisenkom/go-mssqldb
 MsSql: github.com/lunny/godbc
 Oracle: github.com/mattn/go-oci8
 CockroachDB(Postgres): github.com/lib/pq
 ```
 
## 使用教程

> 示例源码https://github.com/zhuxiujia/GoMybatis/tree/master/example

设置好GoPath,用go get 命令下载GoMybatis和对应的数据库驱动
```
go get github.com/zhuxiujia/GoMybatis
go get github.com/go-sql-driver/mysql
```
实际使用mapper
```
import (
	_ "github.com/go-sql-driver/mysql" //导入mysql驱动
	"github.com/zhuxiujia/GoMybatis"
	"fmt"
	"time"
)

//定义xml内容，建议以*Mapper.xml文件存于项目目录中,在编辑xml时就可享受GoLand等IDE渲染和智能提示。生产环境可以使用statikFS把xml文件打包进程序里

var xmlBytes = []byte(`
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE mapper PUBLIC "-//mybatis.org//DTD Mapper 3.0//EN"
"https://raw.githubusercontent.com/zhuxiujia/GoMybatis/master/mybatis-3-mapper.dtd">
<mapper namespace="ActivityMapperImpl">
    <!--SelectAll(result *[]Activity)error-->
    <select id="selectAll">
        select * from biz_activity where delete_flag=1 order by create_time desc
    </select>
</mapper>
`)

type ExampleActivityMapperImpl struct {
     SelectAll  func() ([]Activity, error)
}

func main() {
	var err error
	//Mysql链接格式 用户名:密码@(数据库链接地址:端口)/数据库名称,如root:123456@(***.com:3306)/test
	engine, err := GoMybatis.Open("mysql", "*?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
	   panic(err)
	}
	var exampleActivityMapperImpl ExampleActivityMapperImpl
	
	//加载xml实现逻辑到ExampleActivityMapperImpl
	GoMybatis.WriteMapperPtrByEngine(&exampleActivityMapperImpl, xmlBytes, engine,true)

	//使用mapper
	result, err := exampleActivityMapperImpl.SelectAll(&result)
        if err != nil {
	   panic(err)
	}
	fmt.Println(result)
}
```
## v2019.1.19 新增了 github.com/nylabs/gojee 引擎（改进了原作者源码取消了开头的"."符号,例如.a.b变成 a.b）和 github.com/antonmedv/expr 表达式引擎（改进了原作者源码指针的bug,加入字符串相加操作例如'a'+'b'）
<table border="1">
     <tr>
        <td>表达式引擎</td>
        <td>是否支持指针参数</td>
        <td>执行效率</td>
        <td>表达式功能</td>
    </tr>
    <tr>
         <td>github.com/antonmedv/expr</td>
         <td>支持null和nil和指针</td>
         <td>实测比govaluate快一半以上</td>
         <td>一般</td>
    </tr>
    <tr>
          <td>github.com/nylabs/gojee</td>
          <td>支持null和指针</td>
          <td>实测缓慢-因为每次都有json序列化和反序列化操作</td>
          <td>多</td>
    </tr>
    <tr>
           <td>github.com/Knetic/govaluate</td>
           <td>不支持null和nil和指针</td>
           <td>中等速度</td>
           <td>一般</td>
    </tr>
</table>
为了执行效率 框架默认使用 github.com/antonmedv/expr作为默认选项，你也可以自定义调用GoMybatis.WriteMapper()参数中SqlBuilder的参数自行选择加入ExpressionEngine

## 请及时关注版本，及时升级版本(新的功能，bug修复)
## TODO 期待功能路线（预览特性,有可能会更改）
-针对GoLand和IDEA的xml生成插件，右键一键生成CRUD代码</br>
## 喜欢的老铁欢迎在右上角点下 star 关注和支持我们哈

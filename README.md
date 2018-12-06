# SQL mapper framework for Golang
[![Build Status](https://travis-ci.com/zhuxiujia/GoMybatis.svg?branch=master)](https://travis-ci.com/zhuxiujia/GoMybatis)

![Image text](https://zhuxiujia.github.io/gomybatis.io/assets/vuetify.png)
# 官方网站/文档
https://zhuxiujia.github.io/gomybatis.io/info.html
# 优势
<a href="https://zhuxiujia.github.io/gomybatis.io/info.html">-多种数据库支持</a>,理论上支持mysql和pg的协议以及支持(标准库"database/sql")都支持<br>
<a href="https://zhuxiujia.github.io/gomybatis.io/info.html">-高并发</a>，假设数据库响应时间为0，在6核16Gpc上可框架可以压出 246982Tps,耗时仅仅0.4s<br>
<a href="https://zhuxiujia.github.io/gomybatis.io/info.html">-支持事务</a>，session灵活插拔，兼容过渡期微服务<br>
<a href="https://zhuxiujia.github.io/gomybatis.io/info.html">-动态SQL</a>，在xml中可灵活运用if判断，foreach遍历数组，resultMap,bind等等java框架Mybatis包含的实用功能`<select>,<update>,<insert>,<delete>,<trim>,<if>,<set>,<foreach>,<resultMap>,<bind>,<choose><when><otherwise>`<br>
<a href="https://zhuxiujia.github.io/gomybatis.io/info.html">-无需go generate生成*.go等中间代码</a>，xml读取后可直接写入到自定义的Struct,Func属性中调用函数<br>
### 已支持本地和远程事务,方便处于 单数据库(Mysql,postgresql)-分布式数据库（TiDB,cockroachdb...）过渡期间的微服务
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
	SelectAll         func(result *[]Activity) error
	SelectByCondition func(name string, startTime time.Time, endTime time.Time, page int, size int, result *[]Activity) error `mapperParams:"name,startTime,endTime,page,size"`
	UpdateById        func(session *GoMybatis.Session, arg Activity, result *int64) error //*GoMybatis.Session为事务
	Insert            func(arg Activity, result *int64) error
	CountByCondition  func(name string, startTime time.Time, endTime time.Time, result *int) error                            `mapperParams:"name,startTime,endTime"`
}

func main() {
	var err error
	//Mysql链接格式 用户名:密码@(数据库链接地址:端口)/数据库名称,如root:123456@(***.com:3306)/test
	engine, err := GoMybatis.Open("mysql", "*?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err.Error())
	}
	var exampleActivityMapperImpl ExampleActivityMapperImpl
	
	//挂载xml逻辑到ExampleActivityMapperImpl
	GoMybatis.UseProxyMapperByEngine(&exampleActivityMapperImpl, xmlBytes, engine,true)

	//使用mapper
	var result []Activity
	exampleActivityMapperImpl.SelectAll(&result)

	fmt.Println(result)
}
```

## TODO 期待功能路线（PS 希望同学们踊跃提出创新性功能~~）
-在线编辑SQL支持,可以在系统上线后即时动态修改xml中的SQL逻辑（进行中）</br>
-`<sql><include>` 标签支持（进行中）</br>
-针对于GoLand 的xml生成器,可以一键生成基本的CRUD(待支持..)</br>


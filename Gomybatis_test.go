package GoMybatis

import (
	_ "github.com/go-sql-driver/mysql"

	"fmt"
	"testing"
	"encoding/json"
	"time"
	"github.com/zhuxiujia/GoMybatis/lib/github.com/go-xorm/xorm"
)

var engine *xorm.Engine

func Test_main(t *testing.T) {
	var err error
	engine, err = xorm.NewEngine("mysql", "root:*/test?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err.Error())
	}

	sql := "select count(*) as  C,count(id) as B from biz_activity where delete_flag=1"
	results, err := engine.Query(sql)
	if err != nil {
		panic(err.Error())
	}
	//fmt.Println(string(results[0]["count(*)"]))

	//sql = " select a.*,a.name as Name from biz_activity a where delete_flag=0"
	results, err = engine.Query(sql)
	if err != nil {
		panic(err.Error())
	}

	//var s = make([]map[string]string, 0)
	//for _, v := range results {
	//	var m = make(map[string]string)
	//	for ik, iv := range v {
	//		m[ik] = string(iv)
	//	}
	//	s = append(s, m)
	//}
	//var b, _ = json.Marshal(s)
	//fmt.Println(string(b))

	var activityArray map[string]int

	var n =time.Now()
	for  i:=0;i<1;i++  {
		var e = Unmarshal(results, &activityArray)
		if e != nil {
			panic(e.Error())
		}
	}
	fmt.Println(time.Now().Sub(n).Nanoseconds()/int64(time.Millisecond))
	var bs, _ = json.Marshal(activityArray)
	fmt.Println(string(bs))
}

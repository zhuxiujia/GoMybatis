package GoMybatis

import (
	"testing"
	"time"
)

type TestActivity struct {
	Id         string    `json:"id" gm:"id"`
	Uuid       string    `json:"uuid"`
	Name       string    `json:"name"`
	PcLink     string    `json:"pcLink"`
	H5Link     string    `json:"h5Link"`
	Remark     string    `json:"remark"`
	Version    int       `json:"version" gm:"version"`
	CreateTime time.Time `json:"createTime"`
	DeleteFlag int       `json:"deleteFlag" gm:"logic"`
}

func TestCreateDefaultXml(t *testing.T) {
	var xml = CreateXml("biz_activity", TestActivity{})
	println(string(xml))
}

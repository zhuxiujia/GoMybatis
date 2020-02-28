package GoMybatis

import (
	"testing"
	"time"
)

type TestActivity struct {
	Id         string    `json:"id,omitempty"`
	Uuid       string    `json:"uuid"`
	Name       string    `json:"name"`
	PcLink     string    `json:"pc_link"`
	H5Link     string    `json:"h5_link"`
	Remark     string    `json:"remark"`
	Version    int       `json:"version"`
	CreateTime time.Time `json:"create_time"`
	DeleteFlag int       `json:"delete_flag"`
}

func TestCreateDefaultXml(t *testing.T) {
	var xml = CreateXml("biz_activity", TestActivity{})
	println(string(xml))
}

func TestSnakeString(t *testing.T) {
	println(SnakeString("pcLink"))
}

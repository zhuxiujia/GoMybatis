package GoMybatis

import (
	"github.com/hashicorp/net-rpc-msgpackrpc"
	"net"
	"log"
	"net/rpc"
	"fmt"
	"time"
	"testing"
	"github.com/zhuxiujia/GoMybatis/utils"
)

//注意字段必须是导出
type Params struct {
	Width, Height int
}

type Rect struct{}

func (r *Rect) Area(p *Params, ret *int) error {
	*ret = p.Width + p.Height
	return nil
}

func TestMsgp(t *testing.T) {

	go server()

	var c, e = msgpackrpc.Dial("tcp", "127.0.0.1:1234")
	if e != nil {
		fmt.Println(e)
	}

	time.Sleep(time.Second)

	var total = 100000
	defer utils.CountMethodTps(float64(total), time.Now(), "ZmicroRpcClient")
	ret := 0
	for i := 0; i < total; i++ {
		c.Call("Rect.Area", &Params{50, 100}, &ret)
		fmt.Println(ret)
	}

}
func server() {

	rect := new(Rect)
	//注册rpc服务
	rpc.Register(rect)

	var tcpUrl = "127.0.0.1:1234"

	l, e := net.Listen("tcp", tcpUrl)
	if e != nil {
		log.Fatalf("net rpc.Listen tcp :0: %v", e)
		panic(e)
	}
	for {
		conn, e := l.Accept()
		if e != nil {
			continue
		}
		msgpackrpc.ServeConn(conn)
	}
}

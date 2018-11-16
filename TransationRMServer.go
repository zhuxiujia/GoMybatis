package GoMybatis

import (
	"net/rpc"
	"net"
	"log"
	"github.com/hashicorp/net-rpc-msgpackrpc"
	"fmt"
)

type TransationRMServer struct {
	DefaultTransationManager *DefaultTransationManager
}

func (this TransationRMServer) Msg(arg TransactionReqDTO, result *TransactionRspDTO) {
	fmt.Println(arg)
}

func ServerTransationRM(addr string,driverName, dataSourceName string) {
	transationRMServer := new(TransationRMServer)

	engine, err := Open(driverName, dataSourceName)
	if err != nil {
		panic(err.Error())
	}
	var SessionFactory = SessionFactory{}.New(engine)
	var TransactionFactory = TransactionFactory{}.New(&SessionFactory)
	var manager=DefaultTransationManager{}.New(&SessionFactory, &TransactionFactory)
	transationRMServer.DefaultTransationManager = &manager

	//注册rpc服务
	rpc.Register(transationRMServer)

	var tcpUrl = addr

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

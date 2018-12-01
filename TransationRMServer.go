package GoMybatis

import (
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type TransationRMServer struct {
	DefaultTransationManager *DefaultTransationManager
}

func (this TransationRMServer) Msg(arg TransactionReqDTO, result *TransactionRspDTO) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("[TransationRMServer]work failed:", err)
		}
	}()
	var rsp = this.DefaultTransationManager.DoTransaction(arg)
	*result = rsp
	return nil
}

func ServerTcp(addr string, driverName, dataSourceName string) {
	transationRMServer := new(TransationRMServer)

	engine, err := Open(driverName, dataSourceName)
	if err != nil {
		panic(err.Error())
	}
	var SessionFactory = SessionFactory{}.New(engine)
	var TransactionFactory = TransactionFactory{}.New(&SessionFactory)
	var manager = DefaultTransationManager{}.New(&SessionFactory, &TransactionFactory)
	transationRMServer.DefaultTransationManager = &manager

	//注册rpc服务
	err = rpc.Register(transationRMServer)
	if err != nil {
		panic(err)
	}
	var tcpUrl = addr

	l, e := net.Listen("tcp", tcpUrl)
	if e != nil {
		log.Fatalf("[TransationRMServer]net rpc.Listen tcp :0: %v", e)
		panic(e)
	}
	for {
		conn, e := l.Accept()
		if e != nil {
			continue
		}
		go jsonrpc.ServeConn(conn)
	}
}

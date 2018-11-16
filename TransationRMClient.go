package GoMybatis

import (
	"github.com/hashicorp/net-rpc-msgpackrpc"
	"net/rpc"
)

type TransationRMClient struct {
	Client *rpc.Client
	Addr   string
}

func (this *TransationRMClient) Link(addr string) (*rpc.Client, error) {
	this.Addr = addr
	var client,error= this.autoLink()
	if error!=nil{
		return client,error
	}else {
		this.Client=client
		return client,nil
	}
}
func (this *TransationRMClient) autoLink() (*rpc.Client, error) {
	return msgpackrpc.Dial("tcp", this.Addr)
}

func (this TransationRMClient) Call(arg TransactionReqDTO, result *TransactionRspDTO) {
	this.Client.Call("TransationRMServer.Msg", arg, result)
}

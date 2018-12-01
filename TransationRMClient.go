package GoMybatis

import (
	"errors"
	"net/rpc"
	"net/rpc/jsonrpc"
)

const ConnectError = "connection is shut down"
const CallMethod = "TransationRMServer.Msg"

type RemoteSessionEngine struct {
	SessionEngine
	Client *TransationRMClient
}

func (this RemoteSessionEngine) New(Client *TransationRMClient) RemoteSessionEngine {
	this.Client = Client
	return this
}

func (this *RemoteSessionEngine) NewSession() *Session {
	var TransationSession = TransationRMSession{Client: this.Client}
	var session = Session(&TransationSession)
	return &session
}

type TransationRMClient struct {
	Client    *rpc.Client
	Addr      string
	RetryTime int
}

func (this *TransationRMClient) Link(addr string) (*rpc.Client, error) {
	this.Addr = addr
	var client, error = this.autoLink()
	if error != nil {
		return client, error
	} else {
		this.Client = client
		return client, nil
	}
}
func (this *TransationRMClient) autoLink() (*rpc.Client, error) {
	if this.Client != nil {
		this.Client.Close()
		this.Client = nil
	}
	return jsonrpc.Dial("tcp", this.Addr)
}

func (this *TransationRMClient) Call(arg TransactionReqDTO, result *TransactionRspDTO) error {
	var error error
	if this.Client == nil {
		if this.Addr != "" {
			this.Link(this.Addr)
		} else {
			error = errors.New("[TransationRMClient] link have no addr!")
			return error
		}
	}
	error = this.Client.Call(CallMethod, arg, result)
	if error != nil && error.Error() == ConnectError {
		for i := 0; i < this.RetryTime; i++ {
			this.autoLink()
			error = this.Client.Call(CallMethod, arg, result)
			if error == nil {
				break
			}
		}
	}
	return error
}

package GoMybatis

import (
	"github.com/zhuxiujia/GoMybatis/utils"
	"net/rpc"
	"net/rpc/jsonrpc"
)

const ConnectError = "connection is shut down"
const CallMethod = "TransationRMServer.Msg"

type TransationRMClientConfig struct {
	Addr          string
	RetryTime     int
	TransactionId string
	Status        Transaction_Status
}

type RemoteSessionEngine struct {
	SessionEngine
	Client *TransationRMClient
}

func (it RemoteSessionEngine) New(Client *TransationRMClient) RemoteSessionEngine {
	it.Client = Client
	return it
}

func (it *RemoteSessionEngine) NewSession() Session {
	var TransationSession = TransationRMSession{Client: it.Client}
	var session = Session(&TransationSession)
	return session
}

type TransationRMClient struct {
	Client    *rpc.Client
	Addr      string
	RetryTime int
}

func (it *TransationRMClient) Link(addr string) (*rpc.Client, error) {
	it.Addr = addr
	var client, error = it.autoLink()
	if error != nil {
		return client, error
	} else {
		it.Client = client
		return client, nil
	}
}
func (it *TransationRMClient) autoLink() (*rpc.Client, error) {
	if it.Client != nil {
		it.Client.Close()
		it.Client = nil
	}
	return jsonrpc.Dial("tcp", it.Addr)
}

func (it *TransationRMClient) Call(arg TransactionReqDTO, result *TransactionRspDTO) error {
	var error error
	if it.Client == nil {
		if it.Addr != "" {
			var c, err = it.Link(it.Addr)
			if err != nil {
				return err
			}
			it.Client = c
		} else {
			error = utils.NewError("TransationRMClient", " link have no addr!")
			return error
		}
	}
	error = it.Client.Call(CallMethod, arg, result)
	if error != nil && error.Error() == ConnectError {
		for i := 0; i < it.RetryTime; i++ {
			it.autoLink()
			error = it.Client.Call(CallMethod, arg, result)
			if error == nil {
				break
			}
		}
	}
	return error
}

func (it *TransationRMClient) Close() error {
	if it.Client != nil {
		return it.Client.Close()
	}
	return nil
}

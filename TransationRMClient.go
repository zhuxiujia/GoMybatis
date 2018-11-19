package GoMybatis

import (
	"github.com/hashicorp/net-rpc-msgpackrpc"
	"net/rpc"
	"github.com/kataras/iris/core/errors"
	"github.com/zhuxiujia/GoMybatis/utils"
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
	var TransationSession = RemoteSession{Client: this.Client}
	var session = Session(&TransationSession)
	return &session
}

type RemoteSession struct {
	Session
	SessionId string
	Client *TransationRMClient
}

func (this *RemoteSession)Id() string {
	return this.SessionId
}

func (this *RemoteSession) Query(sqlorArgs string) ([]map[string][]byte, error) {
	var result TransactionRspDTO
	var error = this.Client.Call(TransactionReqDTO{Status: Transaction_Status_NO, TransactionId: utils.CreateUUID(), Sql: sqlorArgs, ActionType: ActionType_Query,}, &result)
	if error == nil && result.Error != "" {
		error = errors.New(result.Error)
	}
	return result.Query, error
}
func (this *RemoteSession) Exec(sqlorArgs string) (Result, error) {
	var result TransactionRspDTO
	var error = this.Client.Call(TransactionReqDTO{Status: Transaction_Status_NO, TransactionId: utils.CreateUUID(), Sql: sqlorArgs, ActionType: ActionType_Exec,}, &result)
	if error == nil && result.Error != "" {
		error = errors.New(result.Error)
	}
	return result.Exec, error
}
func (this *RemoteSession) Rollback() error {
	panic("[RemoteSession] not alow local Rollback()")
	return nil
}
func (this *RemoteSession) Commit() error {
	panic("[RemoteSession] not alow local Commit()")
	return nil
}
func (this *RemoteSession) Begin() error {
	panic("[RemoteSession] not alow local Begin()")
	return nil
}
func (this *RemoteSession) Close() {
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
	return msgpackrpc.Dial("tcp", this.Addr)
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

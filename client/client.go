package client

import (
	"gamecenter.mobi/paicode/wallet"
	"github.com/hyperledger/fabric/peerex"
)

type ClientCore struct{
	Accounts accountManager
	Rpc		 rpcManager
}

func NewClientCore() *ClientCore{
	return &ClientCore{Accounts: accountManager{wallet.CreateSimpleManager("")}}
}

func (c *ClientCore) IsRpcReady() bool{
	return c.Rpc.Rpcbuilder != nil
}


func (c *ClientCore) PrepareRpc(conn peerex.ClientConn){
	c.Rpc.Rpcbuilder = &peerex.RpcBuilder{}
	c.Rpc.Rpcbuilder.Conn = conn
}

func (c *ClientCore) ReleaseRpc(){
	
	if c.Rpc.Rpcbuilder != nil && c.Rpc.Rpcbuilder.Conn.C != nil{
		c.Rpc.Rpcbuilder.Conn.C.Close()	
	}
	
}


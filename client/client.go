package client

import (
	"gamecenter.mobi/paicode/wallet"
	"github.com/hyperledger/fabric/peerex"
)

type ClientCore struct{
	Accounts accountManager
	Rpc		 rpcManager
}

type RpcCore struct{
	Rpc		 rpcManager
}

func NewClientCore(config *peerex.GlobalConfig) *ClientCore{
	
	walletmgr := wallet.CreateSimpleManager(config.GetPeerFS() + "wallet.dat")
	
	return &ClientCore{Accounts: accountManager{walletmgr}}
}

func RpcCoreFromClient(rpc *rpcManager) *RpcCore{
	
	c := new(RpcCore)
	
	c.Rpc.Rpcbuilder = &peerex.RpcBuilder{}
	c.Rpc.Rpcbuilder.Conn = rpc.Rpcbuilder.Conn
	c.Rpc.Rpcbuilder.ChaincodeName = rpc.Rpcbuilder.ChaincodeName
	c.Rpc.PrivKey = rpc.PrivKey
	
	return c
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



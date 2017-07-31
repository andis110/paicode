package client

import (
	"errors"
	
	tx "gamecenter.mobi/paicode/chaincode/transaction"
)

//queryuser: <user id>
func (m* rpcManager) QueryUser(args ...string) ([]byte, error){
	if len(args) != 1{
		return nil, errors.New("Require user id")
	}
	
	m.Rpcbuilder.Function = tx.QueryUser
	return m.Rpcbuilder.Query([]string{args[0]})
		
}

func (m* rpcManager) QueryNode(args ...string) ([]byte, error){
	if len(args) != 0{
		return nil, errors.New("Not require arguments")
	}
	
	m.Rpcbuilder.Function = tx.QueryNode
	return m.Rpcbuilder.Query(nil)
		
}

//queryglobal: <no input>
func (m* rpcManager) QueryGlobal(args ...string) ([]byte, error){
	if len(args) != 0{
		return nil, errors.New("Not require arguments")
	}
	
	m.Rpcbuilder.Function = tx.QueryGlobal
	return m.Rpcbuilder.Query(nil)
		
}

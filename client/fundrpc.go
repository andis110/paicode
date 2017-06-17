package client

import (
	"fmt"
	"errors"
	"strconv"
	
	"gamecenter.mobi/paicode/wallet"
	tx "gamecenter.mobi/paicode/chaincode/transaction"
	txutil "gamecenter.mobi/paicode/transactions"
	"github.com/hyperledger/fabric/peerex"
)

type rpcManager struct{
	PrivKey    wallet.Privkey
	Rpcbuilder *peerex.RpcBuilder
}

const(
	var fundNounceMaxLen int = 256
)


//funding: [to:addr] [amount] <message>
func (m* rpcManager) Fund(args ...string) (string, error){
	if len(args) < 2{
		return errors.New("No required arguments")
	}
	
	b, err := txutil.AddrHelper.VerifyUserId(args[0])
	if !b{
		return err
	}
	
	i, err := strconv.Atoi(args[2])
	if err != nil{
		return err
	}
	
	if i < 0{
		return errors.New("Invalid amount")
	}
	
	fund := &tx.FundTx{FundTxData: tx.FundTxData{args[0], uint(i)}, Invoked: false}
	
	if len(args) == 3{
		if len(args[2]) > fundNounceMaxLen{
			return errors.New("message is too long, should not exceed", fundNounceMaxLen, "chars")
		}
		
		fund.Nounce := []byte(args[2])
	}
	
	rpcargs, err := fund.MakeTransaction(m.PrivKey.K)
	if err != nil{
		return err
	}	
	
	m.Rpcbuilder.Function = tx.UserFund
	return m.Rpcbuilder.Fire(args)
		
}

package main

import (
	"errors"
	"fmt"
	
	"github.com/hyperledger/fabric/core/chaincode/shim"	
	proto "github.com/golang/protobuf/proto"
	
	persistpb "gamecenter.mobi/paicode/protos"
	sec "gamecenter.mobi/paicode/chaincode/security"
	txutil "gamecenter.mobi/paicode/transactions"
	tx "gamecenter.mobi/paicode/chaincode/transaction"

)

func concurrentDuplicateHandle(ud *persistpb.UserData, stub shim.ChaincodeStubInterface) error{
	return nil	
}

func concurrentConsistentHandle(ud *persistpb.UserData, region string) error{
	//regoin checking
	if !sec.Helper.VerifyRegion(ud.ManagedRegion, region){
		return errors.New(fmt.Sprint("User tx is invoked in different region:", region))
	}	
	return nil
}

func (t *PaiChaincode) handleUserFuncs(stub shim.ChaincodeStubInterface, function string, region string, args []string) error{
	
	h, ok := tx.UserTxMap[function]
	if !ok{
		return errors.New(fmt.Sprint("Not a registered function:", function))
	}
	
	cs := txutil.UserTxConsumer{}
	err := cs.ParseArgumentsFirst(args)
	if err != nil{
		return err
	}	
	
	raw, err := stub.GetState(cs.GetUserId())
	if err != nil{
		return err
	}
	
	userdata := &persistpb.UserData{}
	err = proto.Unmarshal(raw, userdata)
	
	if err != nil{
		return err
	}
	
	/*this two step is important for a robust ledger system*/
	err = concurrentConsistentHandle(userdata, region)
	if err != nil{
		return err
	}
	
	err = concurrentDuplicateHandle(userdata, stub)
	if err != nil{
		return err
	}
	
	outuds, err := h.HandleUserTx(userdata, stub, args)
	if err != nil{
		return err
	}
	
	for id, ud := range outuds{
		raw, err := proto.Marshal(ud)
		if err != nil{
			return err
		}
		
		err = stub.PutState(id, raw)
		if err != nil{
			return err
		}
	}
	
	return nil
}


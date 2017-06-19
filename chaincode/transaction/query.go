package transaction

import (
	_ "errors"
	_ "strings"
	
	"github.com/hyperledger/fabric/core/chaincode/shim"
	
)

type queryUserHandler struct{
	
}

type queryGlobalHandler struct{
	
}

func init(){
	QueryMap[QueryUser] = &queryUserHandler{}
	QueryMap[QueryGlobal] = &queryGlobalHandler{}
}

func (_ *querUserHandler) Handle(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){
	
	return nil, nil
} 

func (_ *queryGlobalHandler) Handle(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){
	
	return nil, nil
} 

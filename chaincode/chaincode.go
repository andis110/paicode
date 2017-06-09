package main

import (
	"errors"
	"fmt"
	"sync"
	"strconv"
	"encoding/hex"
	
	"github.com/op/go-logging"
	"github.com/hyperledger/fabric/core/chaincode/shim"	
	proto "github.com/golang/protobuf/proto"
	
	pb "gamecenter.mobi/paicode/protos" 
)


type PaiChaincode struct {
	globalLock sync.RWMutex
	globalSetting *pb.DeploySetting
}

const (

	global_setting_entry string = "global_setting"
	
)

var logger = logging.MustGetLogger("chaincode")

func (t *PaiChaincode) updateCache(stub shim.ChaincodeStubInterface) error{
	t.globalLock.RLock()
	defer t.globalLock.RUnlock()
	
	if pb.DeploySetting == nil{
		t.globalLock.Lock()
		defer t.globalLock.RUnlock()
		
		set, err := stub.GetState(global_setting_entry)
		if err != nil{
			return err
		}
		
		if set == nil{
			return errors.New("FATAL: No global setting found")
		}
		
		t.globalSetting = &pb.DeploySetting{}
		err = proto.Unmarshal(buf, t.globalSetting)
		
		if err != nil{
			return err
		}
		
		logger.Info("Update global setting:", t.globalSetting)	
	}
	
	return nil
}

func (t *PaiChaincode) saveGlobalStatus(stub shim.ChaincodeStubInterface) error{
	t.globalLock.RLock()
	defer t.globalLock.RUnlock()
	
	if pb.DeploySetting == nil{
		return errors.New("FATAL: Invalid cache")
	}	
	
	logger.Info("Save current global setting:", t.globalSetting)	
	
	stat, err := proto.Marshal(t.globalSetting)
	if err != nil{
		return err
	}
	
	return stub.PutState(global_setting_entry, stat)	
}

func (t *PaiChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	return nil, nil
}

func (t *PaiChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	
	if err := t.updateCache(stub); err != nil{
		return nil, err
	}
	
	switch function{
		case "fund":
		case "auth":
	}

	return nil, nil
}

func (t *PaiChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if err := t.updateCache(stub); err != nil{
		return nil, err
	}

	// Get the state from the ledger
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + A + "\"}"
		return nil, errors.New(jsonResp)
	}

	if Avalbytes == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + A + "\"}"
		return nil, errors.New(jsonResp)
	}

	jsonResp := "{\"Name\":\"" + A + "\",\"Amount\":\"" + string(Avalbytes) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return Avalbytes, nil
}

func main() {
	err := shim.Start(new(PaiChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

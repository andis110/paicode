package transaction

import (
	"errors"
	"encoding/base64"
	"crypto/sha256"
	
	proto "github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/timestamp"
	txutil "gamecenter.mobi/paicode/transactions"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

const FundNouncePrefix = "FNC"

type NounceManager struct{
	Tsnow *timestamp.Timestamp
	nouncekey	string
	stub 		shim.ChaincodeStubInterface
}

func (_ *NounceManager) genfundNounce(from string, to string, nounce []byte) (string){
	idbyte := txutil.AddrHelper.DecodeUserid(from)
	if idbyte == nil{
		return ""
	}
	
	shabyte := sha256.Sum256(append(idbyte, nounce...))
	return FundNouncePrefix + base64.StdEncoding.EncodeToString(shabyte[:])
}

//so we get three types return: true and no error indicate we definitely get the exist nounce,
//false and no error indicate we definitely not get the exist nounce or it has been expired,
//false and error indicate we could not know the nounce exist or not and it is on your risk to continue
func (m *NounceManager) CheckfundNounce(stub shim.ChaincodeStubInterface, from string, nounce []byte) (bool, error){
	
	m.nouncekey = m.genfundNounce(from, "", nounce)
	if len(m.nouncekey) == 0{
		return false, errors.New("Could not get func nounce key")
	}
	
	logger.Debug("fund tx nounce:", m.nouncekey)
	m.stub = stub
	
	data, err := stub.GetState(m.nouncekey)
	if err != nil{
		return false, err
	}
	
	if data != nil{
		//nounce can be reused if it has finished for a very long time (nounce_reuse_interval_sec)
		ts := &timestamp.Timestamp{}
		err = proto.Unmarshal(data, ts)
		if err == nil{
			logger.Debug("check nounce's timestamp:", ts, "vs now:", m.Tsnow)
			
			if ts.Seconds + nounce_reuse_interval_sec > m.Tsnow.Seconds{
				return true, nil
			}else{
				return false, nil
			}
		}
		
		return false, err
	}
	
	return false, nil
	
}

func (m *NounceManager) SavefundNounce(){
	if len(m.nouncekey) == 0{
		return
	}
	
	data, err := proto.Marshal(m.Tsnow)
	if err == nil{
		m.stub.PutState(m.nouncekey, data)
	}
	
}


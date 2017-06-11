package transaction

import (
	"fmt"
	"errors"
	"crypto/ecdsa"
	
	"github.com/hyperledger/fabric/core/chaincode/shim"
	
	txutil "gamecenter.mobi/paicode/transactions"
	pb 	   "gamecenter.mobi/paicode/protos"
	persistpb "gamecenter.mobi/paicode/protos"
)

type FundTxData struct{
	To   string
	Amount uint	
}

type FundTx struct{
	FundTxData
	Nounce []byte
	Invoked bool
	InvokedCode uint
}

func (f *FundTxData) fill(v *pb.Funddata) {
	
	f.Amount = uint(v.Pai)
	f.To = v.ToUserId
}

func (f *FundTx) fill(v interface{}) error{
	switch data := v.(type){
		case *pb.UserTxHeader:
		f.Nounce = data.Nounce
		case *pb.Fund:
		switch inndata := data.D.(type){
			case *pb.Fund_Userfund:
			f.FundTxData.fill(inndata.Userfund)
			f.Invoked = false
			case *pb.Fund_InvokeChaincode:
			f.Invoked = true
			f.InvokedCode = uint(inndata.InvokeChaincode)
			default:
			return errors.New(fmt.Sprint("encounter unexpected type in fund field as %T", inndata))
		}
		case *pb.Funddata:
		f.FundTxData.fill(data)
		default:
		return errors.New(fmt.Sprint("encounter unexpected type as %T", data))
	}
	
	return nil
}

func (f *FundTx) Parse(pk *ecdsa.PublicKey, args []string) error{
	
	cs := txutil.UserTxConsumer{PublicKey: pk}

	if len(args) == 3{
		v := &pb.Fund{}
		err := cs.ParseArguments(args, v)
		if err != nil{
			return err
		}
		
		err = f.fill(v)
		if err != nil{
			return err
		}
				
		if f.Invoked {
			return errors.New("A invoked fund tx with not enough arguments")
		}
		
		f.FundTxData.fill(v.D.(*pb.Fund_Userfund).Userfund) 
		
	}else{
		v1 := &pb.Fund{}
		v2 := &pb.Funddata{}
		err := cs.ParseArguments(args, v1, v2)
		if err != nil{
			return err
		}
		
		err = f.fill(v1)
		if err != nil{
			return err
		}

		if !f.Invoked {
			return errors.New("Not a invoked fund tx")
		}
		
		err = f.fill(v2)
		if err != nil{
			return err
		}
	}
	
	return f.fill(cs.HeaderCache)	
}

func (f *FundTx) HandleUserTx(userdata *persistpb.UserData, stub shim.ChaincodeStubInterface, 
	args []string) (outdata []*persistpb.UserData, err error){
	//pk := (*txutil.PublicKey) (userdata.Pk) 
	//cs := txutil.UserTxConsumer{PublicKey: pk.ECDSAPublicKey()}
	
	
	return nil, nil
} 

func (f *FundTx) MakeTransaction(privk *ecdsa.PrivateKey) ([]string, error){
	
	pd := txutil.UserTxProducer{PrivKey: privk, Nounce: f.Nounce}
	fmain := &pb.Fund{}
	fdata := &pb.Funddata{uint32(f.Amount), f.To}
	
	if f.Invoked {
		fmain.D = &pb.Fund_InvokeChaincode{uint32(f.InvokedCode)}		
		return pd.MakeArguments(fmain, fdata)
	}else{
		fmain.D = &pb.Fund_Userfund{fdata}
		return pd.MakeArguments(fmain)
	}
}


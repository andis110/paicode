package transactions

import (
	"fmt"
	"errors"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"math/big"
	"encoding/base64"
	
	"github.com/golang/protobuf/proto"
	pb "gamecenter.mobi/paicode/protos"

)

type FundTxData struct{
	From string
	To   string
	Amount uint	
}

type SignData struct{
	Binding []byte
	signX, signY *big.Int	
}

type FundTx struct{
	FundTxData
	Nounce []byte
	Invoked bool
	InvokedCode uint
}

type FundTxIn struct{
	FundTx
	SignData
}

func (f *FundTxData) fill(v *pb.Funddata) {
	
	f.Amount = uint(v.Pai)
	f.To = v.ToUserId
}

func (f *FundTxIn) fill(v interface{}) error{
	switch data := v.(type){
		case *pb.UserTxHeader:
		f.From = data.FundId
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
		case *pb.Signature:
		f.signX = big.NewInt(0).SetBytes(data.P.X)
		f.signY = big.NewInt(0).SetBytes(data.P.Y)
		case *pb.Funddata:
		f.FundTxData.fill(data)
	}
	
	return nil
}

func (s *SignData) Verify(pk *ecdsa.PublicKey) (bool, error){
	if len(s.Binding) == 0 || s.signX == nil || s.signY == nil{
		return false, errors.New("Not complete data fields")
	}
	
	return ecdsa.Verify(pk, s.Binding, s.signX, s.signY), errors.New("Signature not match")
}

func (f *FundTx) MakeTransaction(privk *ecdsa.PrivateKey) (args []string, err error){
	args = make([]string, 3, 4)
	err = nil
		
	field1 := &pb.UserTxHeader{FundId: f.From, Nounce: f.Nounce}
	field2 := &pb.Fund{}
	field3 := &pb.Signature{}
	field4 := &pb.Funddata{Pai: uint32(f.Amount), ToUserId: f.To}
	fields := append(make([]proto.Message, 0, 4), field1, field2, field3)
	
	if f.Invoked {
		field2.D = &pb.Fund_InvokeChaincode{uint32(f.InvokedCode)}		
		fields = append(fields, field4)
		args = append(args, "")
	}else{
		field2.D = &pb.Fund_Userfund{field4}
	}
	
	hasher := sha256.New()
		
	for i, field := range fields{
		
		if i == 2{
			rx, ry, errx := ecdsa.Sign(rand.Reader, privk, hasher.Sum(nil))
			if errx != nil{
				err = errx
				return
			}			
			field3.P = &pb.ECPoint{rx.Bytes(), ry.Bytes()}
		}
		
		rb, errx := proto.Marshal(field)
		if errx != nil{
			err = errx
			return
		}
		
		if i < 2{
			hasher.Write(rb)
		}
		
		args[i] = base64.StdEncoding.EncodeToString(rb)
	}
	
	return
}


func ParseFundTransaction(args []string) (*FundTxIn, error){
	
	if len(args) < 3 {
		return nil, errors.New(fmt.Sprint("Not enough args, expect at least 3 but only", len(args)))
	}
	
	var pargs []string
	
	ftx := new(FundTxIn)
	hasher := sha256.New()
	
	for i, arg := range args{
		data, err := base64.StdEncoding.DecodeString(arg)
		
		if err != nil{
			return nil, errors.New(fmt.Sprint("base64 decode arg", i, "fail:", err))
		}
		
		var vif proto.Message
		switch i {
			case 0:
				vif = &pb.UserTxHeader{}
			case 1:
				vif = &pb.Fund{}
			case 2:
				if ftx.Invoked && len(pargs) == 3{
					return nil, errors.New("Miss field for invoked transaction")
				}
				vif = &pb.Signature{}
				ftx.Binding = hasher.Sum(nil)
			case 3:
				vif = &pb.Funddata{}
		}
		
		if i < 2{
			hasher.Write(data)
		}
		
		err = proto.Unmarshal(data, vif)
		if err != nil{
			return nil, errors.New(fmt.Sprint("protobuf decode fail", err))
		}
		
		err = ftx.fill(vif)
		if err != nil{
			return nil, errors.New(fmt.Sprint("filling transaction fail", err))
		}
		
		if i == 2 && !ftx.Invoked{
			//done
			break
		}		
	}	
	
	return ftx, nil
}



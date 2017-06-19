package transaction

import (
	"testing"
	"strings"
	"bytes"
	"errors"

	"github.com/hyperledger/fabric/core/chaincode/shim"

	proto "github.com/golang/protobuf/proto"
	"gamecenter.mobi/paicode/wallet"
	pb "gamecenter.mobi/paicode/protos"
	txutil "gamecenter.mobi/paicode/transactions"		
	
)

func compareTest(tx1 *FundTx, tx2 *FundTx) (bool, error){
	if strings.Compare(tx1.To, tx2.To) != 0 ||
		tx1.Amount != tx2.Amount{
			return false, errors.New("funddata not match")
		}
	
	if tx1.Invoked != tx2.Invoked ||
		(tx1.Invoked && tx1.InvokedCode != tx2.InvokedCode){
			return false, errors.New("Invoke fiels not match")
		}
	
	return true, nil
}

func TestTx_UserFund(t *testing.T){

	privk, err := wallet.DefaultWallet.GeneratePrivKey()
	
	if err != nil{
		t.Fatal(err)
	}	
	
	tx1 := &FundTx{FundTxData{"testB", 100}, nil, false, 0}
	
	args, err := tx1.MakeTransaction(privk.K)
	
	if err != nil{
		t.Fatal(err)
	}
	
	if len(args) != 3 {
		t.Fatal("Wrong arg count:", len(args))
	}
	
	t.Log("Output fields", args)
	
	txIn := new(FundTx)
	err = txIn.Parse(&privk.K.PublicKey, args)
	
	if err != nil{
		t.Fatal(err)
	}
	
	t.Log("Output Nounce", txIn.Nounce)
	
	if b, err := compareTest(tx1, txIn); !b{
		t.Fatal(err)
	}
	
	tx2 := &FundTx{FundTxData{"testC", 140}, []byte{44, 44, 44, 44, 44}, true, 13}
	
	args, err = tx2.MakeTransaction(privk.K)
	
	if err != nil{
		t.Fatal(err)
	}
	
	if len(args) != 4 {
		t.Fatal("Wrong arg count:", len(args))
	}
	
	t.Log("Output fields", args)
	
	err = txIn.Parse(&privk.K.PublicKey, args)
	
	if err != nil{
		t.Fatal(err)
	}
	
	t.Log("Output Nounce", txIn.Nounce)
	
	if b, err := compareTest(tx2, txIn); !b{
		t.Fatal(err)
	}
	
	if bytes.Compare(tx2.Nounce, txIn.Nounce) != 0{
		t.Fatal("Nounce not match")
	}	
	
}

func TestFundTx(t *testing.T){

	stub := shim.NewMockStub("DummyTest", nil)	
	
	privk , err := wallet.DefaultWallet.GeneratePrivKey()
	if err != nil{
		t.Fatal(err)
	}
	
	uid := txutil.AddrHelper.GetUserId(&privk.K.PublicKey)
	inpk := privk.GenPublicKeyMsg()
	
	yaprivk , err := wallet.DefaultWallet.GeneratePrivKey()
	if err != nil{
		t.Fatal(err)
	}
	yauid := txutil.AddrHelper.GetUserId(&yaprivk.K.PublicKey)
	
	tx1 := &FundTx{FundTxData{yauid, 100}, []byte{42,42,42}, false, 0}	
	args, err := tx1.MakeTransaction(privk.K)
	if err != nil{
		t.Fatal(err)		
	}	
	
	h := &fundHandler{}
	
	//first tx
	stub.MockTransactionStart("1")
	out, err := h.Handle(uid, &pb.UserData{1000, inpk, nil, "Heaven", nil}, stub, args)
	if err != nil{
		t.Fatal(err)
	}
	stub.MockTransactionEnd("1")
	
	if len(out) != 2{
		t.Fatal("Invalid output")
	}
	
	if v, ok := out[uid]; !ok || v.Pais != 900{
		t.Fatal("Invalid funder")
	} 
	
	if v, ok := out[yauid]; !ok || v.Pais != 100{
		t.Fatal("Invalid accepter")
	}
	
	//duplicated tx, should fail
	stub.MockTransactionStart("1")
	out, err = h.Handle(uid, &pb.UserData{900, inpk, nil, "Heaven", nil}, stub, args)
	if err == nil{
		t.Fatal("Not recognize duplicated tx")
	}
	t.Log(err)
	stub.MockTransactionEnd("1")
	
	data, err := proto.Marshal(&pb.UserData{50, nil, nil, "Heaven", nil})
	if err != nil{
		t.Fatal(err)
	}
	stub.State[yauid] = data
	
	tx2 := &FundTx{FundTxData{yauid, 100}, []byte{42,42,42,42}, false, 0}	
	args, err = tx2.MakeTransaction(privk.K)
	if err != nil{
		t.Fatal(err)		
	}
	
	//another tx, but public key is not match, should fail
	stub.MockTransactionStart("1")
	out, err = h.Handle(uid, &pb.UserData{900, yaprivk.GenPublicKeyMsg(), nil, "Heaven", nil}, stub, args)
	if err == nil{
		t.Fatal("Not recognize wrong publickey")
	}
	t.Log(err)
	stub.MockTransactionEnd("1")
	
	//another tx, but no enough pais, should fail
	stub.MockTransactionStart("1")
	out, err = h.Handle(uid, &pb.UserData{10, inpk, nil, "Heaven", nil}, stub, args)
	if err == nil{
		t.Fatal("Not recognize not enough pais")
	}
	t.Log(err)
	stub.MockTransactionEnd("1")
	
	//another tx
	stub.MockTransactionStart("1")
	out, err = h.Handle(uid, &pb.UserData{900, inpk, nil, "Heaven", nil}, stub, args)
	if err != nil{
		t.Fatal(err)
	}
	stub.MockTransactionEnd("1")
	
	if v, ok := out[uid]; !ok || v.Pais != 800{
		t.Fatal("Invalid funder")
	} 
	
	if v, ok := out[yauid]; !ok || v.Pais != 150{
		t.Fatal("Invalid accepter")
	}else{
		//update accepter data
		data, err := proto.Marshal(v)
		if err != nil{
			t.Fatal(err)
		}
		stub.State[yauid] = data		
	}	
	
	tx3 := &FundTx{FundTxData{yauid, 100}, []byte{42,42,42,42,42}, false, 0}	
	args, err = tx3.MakeTransaction(privk.K)
	if err != nil{
		t.Fatal(err)		
	}	
	
	//yet another tx, make pai become zero
	stub.MockTransactionStart("1")
	out, err = h.Handle(uid, &pb.UserData{100, inpk, nil, "Heaven", nil}, stub, args)
	if err != nil{
		t.Fatal(err)
	}
	stub.MockTransactionEnd("1")
	
	if v, ok := out[uid]; !ok || v.Pais != 0{
		t.Fatal("Invalid funder")
	} 
	
	if v, ok := out[yauid]; !ok || v.Pais != 250{
		t.Fatal("Invalid accepter")
	}	
}


package main

import (
	"testing"
	_ "strings"
	_ "bytes"
	_ "errors"
	"crypto/ecdsa"
	_ "crypto/elliptic"
	_ "crypto/rand"
	
	"github.com/hyperledger/fabric/core/chaincode/shim"
	proto "github.com/golang/protobuf/proto"
	
	paicrypto "gamecenter.mobi/paicode/crypto"
	"gamecenter.mobi/paicode/wallet"
	pb "gamecenter.mobi/paicode/protos"
	tx "gamecenter.mobi/paicode/chaincode/transaction"
	txutil "gamecenter.mobi/paicode/transactions"	
)

const defaultNetCode int32 = 13

func makeInit(stub *shim.MockStub, total int64, preassign map[string]int64) error{
	
	initset := &pb.InitChaincode{&pb.DeploySetting{true, defaultNetCode, total, total}, nil}
	
	for k, v := range preassign{
		initset.PreassignedUser = append(initset.PreassignedUser, &pb.PreassignData{k, v})
	}
	
	arg, err := txutil.EncodeChaincodeTx(initset)
	if err != nil {
		return err
	}
	
	_, err = stub.MockInit("1", "init", []string{arg})
	if err != nil {
		return err
	}
	
	return nil
}

func checkGlobalPai(t *testing.T, stub *shim.MockStub, expect int64 ) {
	buf, ok := stub.State[global_setting_entry]
	if !ok{
		t.Fatal("No global status")
	}
	
	ret := &pb.DeploySetting{}
	err := proto.Unmarshal(buf, ret)
	if err != nil{
		t.Fatal("Unmarshal fail", err)
	}
	
	if !ret.DebugMode || defaultNetCode != ret.NetworkCode || expect != ret.UnassignedPais{
		t.Fatal("Not correct global setting", ret)
	}
}

func checkUser(t *testing.T, stub *shim.MockStub, uid string, expect int64){
	buf, ok := stub.State[uid]
	if !ok{
		t.Fatal("No user", uid)
	}
	
	ret := &pb.UserData{}
	err := proto.Unmarshal(buf, ret)
	if err != nil{
		t.Fatal("Unmarshal fail", err)
	}
	
	if ret.Pais != expect{
		t.Fatal("Not correct pais for user", uid, ret.Pais, expect)
	}
}

func TestPaichaincode_Init(t *testing.T) {
	pcc := new(PaiChaincode)
	stub := shim.NewMockStub("PaicodeTest", pcc)

	err := makeInit(stub, 100000, map[string]int64{})
	if err != nil{
		t.Fatal(err)
	}

	checkGlobalPai(t, stub, 100000)
}

func TestPaichaincode_InitPreassign(t *testing.T) {
	pcc := new(PaiChaincode)
	stub := shim.NewMockStub("PaicodeTest", pcc)

	err := makeInit(stub, 100000, map[string]int64{"dummy1": 50000, "dummy2": 10})
	if err != nil{
		t.Fatal(err)
	}

	checkGlobalPai(t, stub, 49990)
	checkUser(t, stub, "dummy1", 50000)
	checkUser(t, stub, "dummy2", 10)
}

type privKey struct{
	k 			*ecdsa.PrivateKey
	underlyingK	*paicrypto.ECDSAPriv
} 

func producePrivk(count int) (ret []*wallet.Privkey){
	
	ret = make([]*wallet.Privkey, count)
	for i, _ := range ret{
		k , _ := wallet.DefaultWallet.GeneratePrivKey()
		if k == nil{
			ret = nil
			return
		}
		
		ret[i] = k
	}
	
	return
} 

func confirmUser(t *testing.T, stub *shim.MockStub, privk *wallet.Privkey){
	pd := &txutil.UserTxProducer{PrivKey: privk.K}
	
	args, err := pd.MakeArguments(&pb.RegPublicKey{privk.GenPublicKeyMsg()})
	
	if err != nil{
		t.Fatal(err)		
	}
	
	_, err = stub.MockInvoke("confirmTest", tx.UserRegPublicKey, args)
	if err != nil{
		t.Fatal(err)
	}
}

func testFailConfirmUser(t *testing.T, stub *shim.MockStub, privk *wallet.Privkey, wrong_privk *wallet.Privkey){
	pd := &txutil.UserTxProducer{PrivKey: privk.K}
	
	args, err := pd.MakeArguments(&pb.RegPublicKey{wrong_privk.GenPublicKeyMsg()})
	
	if err != nil{
		t.Fatal(err)		
	}
	
	_, err = stub.MockInvoke("confirmTest", tx.UserRegPublicKey, args)
	if err == nil{
		t.Fatal("Not fail wrong user reg")
	}
	
	t.Log(err)
}

func TestPaichaincode_FundTx(t *testing.T) {
	pcc := new(PaiChaincode)
	stub := shim.NewMockStub("PaicodeTest", pcc)

	keys := producePrivk(3)
	if len(keys) < 3{
		t.Fatal("Produce keys fail")
	}

	ids := [3]string{
		txutil.AddrHelper.GetUserId(&keys[0].K.PublicKey),
		txutil.AddrHelper.GetUserId(&keys[1].K.PublicKey),
		txutil.AddrHelper.GetUserId(&keys[2].K.PublicKey)}

	err := makeInit(stub, 100000, map[string]int64{
			ids[0]: 50000, 
			ids[1]: 10,
			ids[2]: 99})
	
	if err != nil{
		t.Fatal(err)
	}

	confirmUser(t, stub, keys[0])
	confirmUser(t, stub, keys[2])
	testFailConfirmUser(t, stub,  keys[1],  keys[2])

}



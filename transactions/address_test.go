package transactions

import (
	"testing"
	"strings"
	"math/big"
	"crypto/rand"
	
	paicrypto "gamecenter.mobi/paicode/crypto"
)

const (
	testPrefix int = 1
)

func testIdA(t *testing.T, pk1 *PublicKey, pk2 *PublicKey){

	id1 := pk1.GetUserId(testPrefix)
	id2 := pk2.GetUserId(testPrefix)
	
	if len(id1) == 0 || len(id2) == 0{
		t.Fatal("Invalid id", id1, "--", id2)
	}
	
	if strings.Compare(id1, id2) != 0{
		t.Fatal("Not identify id", id1, "--", id2)
	}
	
	t.Log(id1)
	
}

func testIdB(t *testing.T, pk1 *PublicKey, pk2 *PublicKey){

	id1 := pk1.GetUserId(testPrefix)
	id2 := pk2.GetUserId(testPrefix)
	
	if len(id1) == 0 || len(id2) == 0{
		t.Fatal("Invalid id", id1, "--", id2)
	}
	
	if strings.Compare(id1, id2) == 0{
		t.Fatal("Not unique id", id1, "--", id2)
	}
	
	t.Log(id1, "--", id2)
	
}

func TestDump_Userid(t *testing.T){
	
	rb := make([]byte, 32)
	_, err := rand.Read(rb)
	if err != nil{
		t.Skip("rand make 256bit bytes fail", err)
	}	
	
	var one, sed1, sed2 *big.Int
	one = big.NewInt(1) 
	sed1 = new(big.Int)
	sed2 = new(big.Int)
	
	sed1.SetBytes(rb)
	sed2.Add(sed1, one)
	
	publick1 := NewPublicKeyFromPriv(&paicrypto.ECDSAPriv{paicrypto.ECP256_FIPS186, sed1})
	publick2 := NewPublicKeyFromPriv(&paicrypto.ECDSAPriv{paicrypto.ECP256_FIPS186, sed2})

	if publick1 == nil || publick2 == nil{
		t.Fatal("Invalid public key")
	}
	
	_, err = publick1.GetUserHash()
	if err != nil{
		t.Fatal("fail hash:", err)
	}
	
	testIdA(t, publick1, publick1)
	testIdB(t, publick1, publick2)
	
	for i := 0; i < 5000; i++ {
		sed2 = sed2.Add(sed2, one)
		publick1 = publick2
		testIdA(t, publick1, publick2)
		publick2 = NewPublicKeyFromPriv(&paicrypto.ECDSAPriv{paicrypto.ECP256_FIPS186, sed2})
		testIdB(t, publick1, publick2)
	}
	
	
}

func TestVerify_Userid(t *testing.T){
	
	rb := make([]byte, 32)
	_, err := rand.Read(rb)
	if err != nil{
		t.Skip("rand make 256bit bytes fail", err)
	}	
	
	sed := new(big.Int)	
	sed.SetBytes(rb)
	
	publick := NewPublicKeyFromPriv(&paicrypto.ECDSAPriv{paicrypto.ECP256_FIPS186, sed})
	
	uid := publick.GetUserId(1)
	
	if b, err := VerifyUserId(uid, 1); !b{
		t.Fatal("verify fail 1", err)
	}
	
	if b, err := VerifyUserId(uid, 0); b{
		t.Fatal("verify error 2")		
	}else{
		t.Log("verfiy ret", err)
	}
}

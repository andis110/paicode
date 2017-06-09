package transactions

import (
	"testing"
	"crypto/rand"
	"crypto/ecdsa"
	"crypto/elliptic"
	
	"github.com/golang/protobuf/proto"
	pb "gamecenter.mobi/paicode/protos"
	paicrypto "gamecenter.mobi/paicode/crypto"
)

func TestPublicKey(t *testing.T){

	prvkeystd, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil{
		t.Skip("Skip for ecdsa lib fail:", err)
	}
	
	rb := make([]byte, 32)
	_, err = rand.Read(rb)
	if err != nil{
		t.Skip("rand make 256bit bytes fail", err)
	}	
	
	publick := NewPublicKeyFromPriv(&paicrypto.ECDSAPriv{paicrypto.ECP256_FIPS186, prvkeystd.D})
	
	if publick == nil{
		t.Fatal("Generate public key from NewPublicKeyFromPriv fail")
	}		
	
	sx, sy, err := ecdsa.Sign(rand.Reader, prvkeystd, rb)
	if err != nil{
		t.Fatal(err)
	}
	
	if !ecdsa.Verify(&publick.PublicKey, rb, sx, sy) {
		t.Fatal("verify signature with generated publickey fail")
	}
	
}

func TestPublicKey_DumpPb(t *testing.T){

	prvkeystd, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil{
		t.Skip("Skip for ecdsa lib fail:", err)
	}
	
	rb := make([]byte, 32)
	_, err = rand.Read(rb)
	if err != nil{
		t.Skip("rand make 256bit bytes fail", err)
	}	
	
	publick := NewPublicKeyFromPriv(&paicrypto.ECDSAPriv{paicrypto.ECP256_FIPS186, prvkeystd.D})
	
	if publick == nil{
		t.Fatal("Generate public key from NewPublicKeyFromPriv fail")
	}	
	
	pbdump := publick.MakePb()
	
	if pbdump == nil {
		t.Fatal("Make protobuf message fail")
	}
	
	msgbyte, err := proto.Marshal(pbdump)
	if err != nil{
		t.Fatal("Marshal protobuf fail", err)
	}
	
	pbrcv := new(pb.PublicKey)
	
	err = proto.Unmarshal(msgbyte, pbrcv)
	if err != nil{
		t.Fatal("Unmarshal protobuf fail", err)
	}
	
	publick2 := NewPublicKey(pbrcv)
	
	if publick2 == nil{
		t.Fatal("Generate public key from NewPublicKey fail")
	}
	
	sx, sy, err := ecdsa.Sign(rand.Reader, prvkeystd, rb)
	if err != nil{
		t.Fatal(err)
	}
	
	if !ecdsa.Verify(&publick2.PublicKey, rb, sx, sy) {
		t.Fatal("verify signature with dump publickey fail")
	}	
	
}

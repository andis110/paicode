package transactions

import (
	"testing"
	"errors"
	"bytes"
	"strings"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
)

func compareTest(tx1 *FundTx, tx2 *FundTx) (bool, error){
	if strings.Compare(tx1.From, tx2.From) != 0 ||
		strings.Compare(tx1.To, tx2.To) != 0 ||
		tx1.Amount != tx2.Amount{
			return false, errors.New("funddata not match")
		}
	
	if tx1.Invoked != tx2.Invoked ||
		(tx1.Invoked && tx1.InvokedCode != tx2.InvokedCode){
			return false, errors.New("Invoke fiels not match")
		}
	
	if bytes.Compare(tx1.Nounce, tx2.Nounce) != 0{
		return false, errors.New("Nounce not match")
	}
	
	return true, nil
}

func TestTx_UserFund(t *testing.T){

	privk, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	
	if err != nil{
		t.Skip("Skip for ecdsa lib fail:", err)
	}	
	
	tx1 := &FundTx{FundTxData{"testA", "testB", 100}, []byte{42, 42, 42, 42}, false, 0}
	
	args, err := tx1.MakeTransaction(privk)
	
	if err != nil{
		t.Fatal(err)
	}
	
	if len(args) != 3 {
		t.Fatal("Wrong arg count:", len(args))
	}
	
	t.Log("Output fields", args)
	
	txIn1, err := ParseFundTransaction(args)
	
	if err != nil{
		t.Fatal(err)
	}
	
	if b, err := txIn1.Verify(&privk.PublicKey); !b{
		t.Fatal(err)
	}
	
	if b, err := compareTest(tx1, &txIn1.FundTx); !b{
		t.Fatal(err)
	}
	
	tx2 := &FundTx{FundTxData{"testB", "testC", 140}, []byte{44, 44, 44, 44, 44}, true, 13}
	
	args, err = tx2.MakeTransaction(privk)
	
	if err != nil{
		t.Fatal(err)
	}
	
	if len(args) != 4 {
		t.Fatal("Wrong arg count:", len(args))
	}
	
	t.Log("Output fields", args)
	
	txIn2, err := ParseFundTransaction(args)
	
	if err != nil{
		t.Fatal(err)
	}
	
	if b, err := txIn2.Verify(&privk.PublicKey); !b{
		t.Fatal(err)
	}
	
	if b, err := compareTest(tx2, &txIn2.FundTx); !b{
		t.Fatal(err)
	}	
	
}
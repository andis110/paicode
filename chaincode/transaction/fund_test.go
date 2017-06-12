package transaction

import (
	"testing"
	"strings"
	"bytes"
	"errors"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	
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

	privk, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	
	if err != nil{
		t.Skip("Skip for ecdsa lib fail:", err)
	}	
	
	tx1 := &FundTx{FundTxData{"testB", 100}, nil, false, 0}
	
	args, err := tx1.MakeTransaction(privk)
	
	if err != nil{
		t.Fatal(err)
	}
	
	if len(args) != 3 {
		t.Fatal("Wrong arg count:", len(args))
	}
	
	t.Log("Output fields", args)
	
	txIn := new(FundTx)
	err = txIn.Parse(&privk.PublicKey, args)
	
	if err != nil{
		t.Fatal(err)
	}
	
	t.Log("Output Nounce", txIn.Nounce)
	
	if b, err := compareTest(tx1, txIn); !b{
		t.Fatal(err)
	}
	
	tx2 := &FundTx{FundTxData{"testC", 140}, []byte{44, 44, 44, 44, 44}, true, 13}
	
	args, err = tx2.MakeTransaction(privk)
	
	if err != nil{
		t.Fatal(err)
	}
	
	if len(args) != 4 {
		t.Fatal("Wrong arg count:", len(args))
	}
	
	t.Log("Output fields", args)
	
	err = txIn.Parse(&privk.PublicKey, args)
	
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

func TestTx_HandlingFund(t *testing.T){
	
}

}

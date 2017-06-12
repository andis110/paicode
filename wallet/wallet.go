package wallet

import (
	"crypto/rand"
	"crypto/ecdsa"

	pb		  "gamecenter.mobi/paicode/protos"
	paicrypto "gamecenter.mobi/paicode/crypto"
	txutil	  "gamecenter.mobi/paicode/transactions"
)

type Wallet struct{
	useCurve int
}

type Privkey struct{
	K 	*ecdsa.PrivateKey
	underlyingKey *paicrypto.ECDSAPriv
}

func (k Privkey) GenPublicKeyMsg() *pb.PublicKey{
	ret, err := txutil.MakePbFromPrivKey(k.underlyingKey)
	if err != nil{
		return nil
	}
	
	return ret
}

var DefaultWallet = Wallet{
	useCurve: paicrypto.ECP256_FIPS186}

func (w *Wallet) GeneratePrivKey() (*Privkey, error){
	
	curve, err := paicrypto.GetEC(w.useCurve)
	if err != nil{
		return nil, err
	}
	
	ecprivk, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil{
		return nil, err
	}
	
	return &Privkey{ecprivk, &paicrypto.ECDSAPriv{w.useCurve, ecprivk.D}}, nil
}

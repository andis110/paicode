package wallet

import (
	"crypto/rand"
	"crypto/ecdsa"

	paicrypto "gamecenter.mobi/paicode/crypto"
)

type Wallet struct{
	useCurve int
}

var DefaultWallet = Wallet{
	useCurve: paicrypto.ECP256_FIPS186}

func (w *Wallet) GeneratePrivKey() (error, *ecdsa.PrivateKey, *paicrypto.ECDSAPriv){
	
	curve, err := paicrypto.GetEC(w.useCurve)
	if err != nil{
		return err, nil, nil
	}
	
	ecprivk, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil{
		return err, nil, nil
	}
	
	return nil, ecprivk, &paicrypto.ECDSAPriv{w.useCurve, ecprivk.D}
}

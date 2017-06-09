package transactions

import (
	"crypto/ecdsa"
	"math/big"
	
	pb "gamecenter.mobi/paicode/protos"
	paicrypto "gamecenter.mobi/paicode/crypto"

)

type PublicKey struct{
	ecdsa.PublicKey
	CurveType int32
}

func NewPublicKey(pk *pb.PublicKey) *PublicKey{
	curve, err := paicrypto.GetEC(int(pk.Curvetype))
	if err != nil{
		return nil
	}
	
	return &PublicKey{ecdsa.PublicKey{curve, big.NewInt(0).SetBytes(pk.P.GetX()),
		big.NewInt(0).SetBytes(pk.P.GetY())}, pk.Curvetype}	
}

func NewPublicKeyFromPriv(priv *paicrypto.ECDSAPriv) *PublicKey{
	pk, err := priv.Apply()
	
	if err != nil{
		return nil
	}
	
	return &PublicKey{pk.PublicKey, int32(priv.CurveType)}
}

func (pk *PublicKey) MakePb() *pb.PublicKey{
	
	return &pb.PublicKey{pk.CurveType, &pb.ECPoint{pk.X.Bytes(), pk.Y.Bytes()}}
	
}


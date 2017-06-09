package transactions

import (
	"errors"
	"bytes"
	"math/big"
	"encoding/asn1"
	"encoding/base64"
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
)

func (pk *PublicKey) GetUserHash() ([]byte, error) {
	
	type StandardPk struct{
		X *big.Int
		Y *big.Int
		CurveType int "tag:10"
	}
	
	stdpk := StandardPk{X: pk.X, Y: pk.Y, CurveType: int(pk.CurveType)}
	
	rawbytes, err := asn1.Marshal(stdpk)
	
	if err != nil{
		return nil, err
	}
	
	hash1pass := sha256.Sum256(rawbytes)
	if len(hash1pass) != sha256.Size{
		return nil, errors.New("Wrong sha256 hashing")
	}
	
	rmd160h := ripemd160.New();
	if nn, err := rmd160h.Write(hash1pass[:]); nn != len(hash1pass) || err != nil{
		return nil, errors.New("Wrong ripemd write")
	}
	
	hash2pass := rmd160h.Sum([]byte{})
	
	if len(hash2pass) != ripemd160.Size{
		return nil, errors.New("Wrong ripemd160 hashing")
	}
	
	return hash2pass, nil
}

const (
	AddressFullByteSize = 25
	AddressPartByteSize = 21
	AddressVerifyCodeSize = 4
)

func getCheckSum(rb []byte) ([AddressVerifyCodeSize]byte, error){
	
	var ret [AddressVerifyCodeSize]byte
	
	hash1pass := sha256.Sum256(rb)
	if len(hash1pass) != sha256.Size{
		return ret, errors.New("Wrong sha256 hashing 1pass")
	}	
	
	hash2pass := sha256.Sum256(hash1pass[:])
	if len(hash2pass) != sha256.Size{
		return ret, errors.New("Wrong sha256 hashing 2pass")
	}
	
	return [AddressVerifyCodeSize]byte{hash2pass[0], hash2pass[1], hash2pass[2], hash2pass[3]}, nil
}

func VerifyUserId(id string, prefix int) (bool, error){
	data, err := base64.StdEncoding.DecodeString(id)
	
	if uint8(prefix % 256) != uint8(data[0]){
		return false, errors.New("Different prefix")
	} 
	
	if err != nil {
		return false, errors.New("Wrong base64 decoding")
	}
	
	if len(data) != AddressFullByteSize{
		return false, errors.New("Wrong bytes size")
	}
	
	ck, err := getCheckSum(data[:AddressPartByteSize])
	if err != nil{
		return false, errors.New("Get checksum fail")
	}
	
	return bytes.Equal(ck[:], data[AddressPartByteSize:]), errors.New("checksum not equal")
}

func GenUserId(prefix int, rb []byte) string{

	fullbytes := make([]byte, 1, AddressFullByteSize)
	fullbytes[0] = uint8(prefix % 256)	
	fullbytes = append(fullbytes, rb...)
	
	if len(fullbytes) != AddressPartByteSize {
		return ""
	}
	
	ck, err := getCheckSum(fullbytes)	
	if err != nil{
		return ""
	}
	
	return base64.StdEncoding.EncodeToString(append(fullbytes, ck[:]...))
	
}

func (pk *PublicKey) GetUserId(prefix int) string{
		
	hashbytes, err := pk.GetUserHash()
	if err != nil{
		return ""
	}
	
	return GenUserId(prefix, hashbytes)
}


	
	




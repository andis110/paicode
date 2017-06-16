package client

import (
	"fmt"
	"errors"
	
	"gamecenter.mobi/paicode/wallet"
)

type accountManager struct{
	KeyMgr *wallet.KeyManager
}

//generate privatekey: <remark>
func (m* accountManager) GenPrivkey(args []string) error{
	if len(args) > 1{
		return errors.New(fmt.Sprint("Could not recognize", args[1:]))
	}
	
	var remark string		
	if len(args) == 0{
		remark = RandStringRunes(16)
	}else{
		remark = args[0]
	}
	
	k, err := wallet.DefaultWallet.GeneratePrivKey()
	if err != nil{
		return err
	}
	
	m.KeyMgr.AddPrivKey(remark, k)
	
	return nil
		
}

//dump privatekey from [remark]
func (m* accountManager) DumpPrivkey(args []string) (string, error){
	if len(args) != 1{
		return errors.New("Invalid remark")
	}
	
	k, err := m.KeyMgr.LoadPrivKey(args[0])
	if err != nil{
		return err
	}
	
	return k.DumpPrivkey()
}

//get address from [remark]
func (m* accountManager) GetAddress(args []string) (string, error){
	if len(args) != 1{
		return errors.New("Invalid remark")
	}
	
	k, err := m.KeyMgr.LoadPrivKey(args[0])
	if err != nil{
		return err
	}

	
	
	return "", nil
}

//list all keys in manager with remark and address
func (m* accountManager) ListKeyData(args []string) [][2]string{
	if len(args) != 0{
		return errors.New("No argument required")
	}
		
	return nil
}

//[import string], <remark>
func (m* accountManager) ImportPrivkey(args []string) error{
	
	if len(args) == 0{
		return errors.New("Need import string")
	}	
	
	if len(args) > 2{
		return errors.New(fmt.Sprint("Could not recognize", args[2:]))
	}
	
	var remark string		
	if len(args) == 1{
		remark = RandStringRunes(16)
	}else{
		remark = args[1]
	}	
	
	k, err := m.KeyMgr.ImportPrivKey(args[0])
	if err != nil{
		return err
	}
	
	m.KeyMgr.AddPrivKey(remark, k)
	
	return nil
}


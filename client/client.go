package client

import (
	"gamecenter.mobi/paicode/wallet"
)

type ClientCore struct{
	Accounts accountManager
}

func NewClientCore() *ClientCore{
	return &ClientCore{Accounts: accountManager{wallet.CreateSimpleManager("")}}
}

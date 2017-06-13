package wallet

type KeyManager interface {
	
	AddPrivKey(remark string, privk *Privkey)
	LoadPrivKey(remark string) (*Privkey, error)
	ListAll() (map[string]*Privkey, error)
	Load() error
	Persist() error
}
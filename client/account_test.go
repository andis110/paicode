package client

import (
    "testing"
    "strings"
    "gamecenter.mobi/paicode/wallet"
)

func TestMgrPrivk(t *testing.T) {
	mgr := accountManager{wallet.CreateSimpleManager("")}
	
	_, err := mgr.GenPrivkey("test1")
	if err != nil{
		t.Fatal(err)
	}
	
	dstr, err := mgr.DumpPrivkey("test1")
	if err != nil{
		t.Fatal(err)
	}
	
	_, err = mgr.DumpPrivkey("test2")
	if err == nil{
		t.Fatal("Dump unexist key")
	}
	
	err = mgr.ImportPrivkey(dstr)
	if err != nil{
		t.Fatal(err)
	}
	
	err = mgr.ImportPrivkey(dstr, "test2")
	if err != nil{
		t.Fatal(err)
	}
	
	dstr2, err := mgr.DumpPrivkey("test2")
	if err != nil{
		t.Fatal(err)
	}
	
	if strings.Compare(dstr, dstr2) != 0{
		t.Fatal("Dumped key not identical")
	}
	
	list := mgr.ListKeyData()
	if len(list) != 3{
		t.Fatal("Not expected count for keys")
	}
	
	t.Log(list)
	
	_, err = mgr.GetAddress("test2")
	if err != nil{
		t.Fatal(err)
	}
	
	_, err = mgr.GetAddress("test3")
	if err == nil{
		t.Fatal("Get unexist address")
	}
		
}


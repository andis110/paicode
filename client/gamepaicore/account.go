package main

import (
	
	"fmt"
	"net/http"
	"github.com/gocraft/web"
	"encoding/json"
	_ "gamecenter.mobi/paicode/client"
)

type accountData struct{
	Status string
	Data   interface{}
}

func (s *AccountREST) QueryAcc(rw web.ResponseWriter, req *web.Request){

	id := req.PathParams["id"]
	if id == "" {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("{\"Status\":\"Account id not found\"}"))//just write a raw json
		return
	}
	
	addr, err := defClient.Accounts.GetAddress(id)
	if err != nil{
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte(fmt.Sprintf("{\"Status\":\"Account not exist: %s\"}", err)))
		return
	}
	
	encoder := json.NewEncoder(rw)
	rw.WriteHeader(http.StatusOK)
	encoder.Encode(accountData{"ok", addr})
	
}

func (s *AccountREST) ListAcc(rw web.ResponseWriter, req *web.Request){
	
	ret := defClient.Accounts.ListKeyData()
	
	encoder := json.NewEncoder(rw)
	rw.WriteHeader(http.StatusOK)
	encoder.Encode(accountData{"ok", ret})
}

func (s *AccountREST) NewAcc(rw web.ResponseWriter, req *web.Request){
	
	err := req.ParseForm()
	if err != nil || len(req.Form) == 0{
		rw.WriteHeader(http.StatusNotAcceptable)
		rw.Write([]byte("{\"Status\":\"Wrong form or not expected content (application/x-www-urlencoded)\"}"))
		return		
	}
		
	accountid := req.Form["id"]	
	if len(accountid) == 0{
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("{\"Status\":\"Must provide account id\"}"))
		return				
	}
	
	prvkstr := req.Form["privatekey"]
		
	if len(prvkstr) != 0{
		//import
		_, err = defClient.Accounts.ImportPrivkey(prvkstr[0], accountid[0])
	}else{
		//generate
		_, err = defClient.Accounts.GenPrivkey(accountid[0])
	}
	
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("{\"Status\":\"ok\"}"))
	s.shouldPersist = true
}

func (s *AccountREST) DumpAcc(rw web.ResponseWriter, req *web.Request){
	
	//TODO
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("{\"Status\":\"ok\"}"))	
}


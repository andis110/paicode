package main 

import (
	"fmt"
	"github.com/gocraft/web"

	clicore "gamecenter.mobi/paicode/client"
	

)

type GamepaiREST struct{
	
}

type AccountREST struct{
	shouldPersist bool
}

type RpcREST struct{
	workCore *clicore.RpcCore 
}

func (s *GamepaiREST) SetResponseType(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	rw.Header().Set("Content-Type", "application/json")

	// Enable CORS (default option handler will handle OPTION and set Access-Control-Allow-Method properly)
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Headers", "accept, content-type")

	next(rw, req)
}

func (s *AccountREST) PersistAccount(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {

	next(rw, req)
	
	if s.shouldPersist {
		err := defClient.Accounts.KeyMgr.Persist()
		
		if err != nil{
			panic(fmt.Sprintln("Persist fail", err))
		}
		
	} 
}


func (s *RpcREST) LoadRPC(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {

	s.workCore = clicore.RpcCoreFromClient(&defClient.Rpc)
	kid := req.PathParams["id"]
	if kid == "" {
		panic("Must specific id")
	}
	
	key, err := defClient.Accounts.KeyMgr.LoadPrivKey(kid)
	if err != nil{
		panic(fmt.Sprintf("No corresponding privkey for %s", kid))
	}	
	
	s.workCore.Rpc.PrivKey = key	
	next(rw, req)
}

func buildRouter() *web.Router {
	
	router := web.New(GamepaiREST{})

	// Add middleware
	router.Middleware((*GamepaiREST).SetResponseType)

	accRouter := router.Subrouter(AccountREST{false}, "/account")
	accRouter.Middleware((*AccountREST).PersistAccount)
	accRouter.Post("/", (*AccountREST).NewAcc)
	accRouter.Get("/", (*AccountREST).ListAcc)
	accRouter.Get("/:id", (*AccountREST).QueryAcc)	
	accRouter.Get("/dump/:id", (*AccountREST).DumpAcc)
	
	rpcRouter := router.Subrouter(RpcREST{}, "/rpc")
	rpcRouter.Middleware((*RpcREST).LoadRPC)
	rpcRouter.Get("/:id", (*RpcREST).Query) //query user	
	rpcRouter.Post("/:id/registar", (*RpcREST).Registar)
	rpcRouter.Post("/:id/fund", (*RpcREST).Fund)

	return router
}


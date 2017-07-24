package main 

import (
	"os"
	"fmt"
	"net/http"
	_ "github.com/gocraft/web"
	"github.com/spf13/cobra"
	
	clicore "gamecenter.mobi/paicode/client"
	"github.com/hyperledger/fabric/peerex"
)

const defPaicodeName string = "50637ebc88e9c0f2ea9d240784b491c4fde8ebd177a95fbc2f087312111affef1898fea4c267ff1084db244de6c6860f4367b700659d44b7b47fabda27347c23"

var defClient *clicore.ClientCore 

var mainCmd = &cobra.Command{
	Use: "gamepai",
	
	PreRunE: func(cmd *cobra.Command, args []string) error {

		config := &peerex.GlobalConfig{}
		err := config.InitGlobal()
		
		if err != nil{
			return err
		}
		
		err = os.MkdirAll(config.GetPeerFS(), 0777)
		if err != nil{
			return err
		}	
		
		defClient = clicore.NewClientCore(config)

		if !offlinemode {
			
			conn := peerex.ClientConn{}
			err := conn.Dialdefault()
			if err != nil{
				return err
			}			
			
			defClient.PrepareRpc(conn)
			defClient.Rpc.Rpcbuilder.ChaincodeName = defPaicodeName
			restLogger.Infof("Start rpc, chaincode is %s", defClient.Rpc.Rpcbuilder.ChaincodeName)
				
		}else{
			restLogger.Info("Run under off-line mode")
		}
		
		return nil

	},
	
	Run: func(cmd *cobra.Command, args []string){
		
		var svraddr string
		if len(args) > 1{
			svraddr = args[0]
		}else{
			svraddr = "localhost:7280"
		}
		
		defClient.Accounts.KeyMgr.Load()
		defer defClient.Accounts.KeyMgr.Persist()			
		
		// Initialize the REST service object
		restLogger.Infof("Initializing the REST service on %s", svraddr)
	
		router := buildRouter()
		err := http.ListenAndServe(svraddr, router)
		if err != nil {
			restLogger.Errorf("ListenAndServe: %s", err)
		}
		
		if defClient.IsRpcReady(){
			defClient.ReleaseRpc()
		}
	},	
}

var exitCmd = &cobra.Command{
	Use: "exit",
	Run: func(cmd *cobra.Command, args []string){
		//TODO, call exit API?
	},
}

var restLogger = peerex.InitLogger("gamepaiREST")
var debugmode bool = false
var offlinemode bool = false


func main() {
	
	mainCmd.Flags().BoolVar(&debugmode, "debug", false, "run http server with debug output")
	mainCmd.Flags().BoolVar(&offlinemode, "offline", false, "not communicate with other peers")
	
	mainCmd.AddCommand(exitCmd)	

	err := mainCmd.Execute()
	if err != nil{
		fmt.Println("Command handler error:", err)
		os.Exit(1)		
	}

}


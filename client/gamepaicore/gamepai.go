package main // import "gamecenter.mobi/paicode/client/gamepaicore"

import (
	"os"
	"fmt"
	"net/http"
	_ "github.com/gocraft/web"
	"github.com/spf13/cobra"
	
	clicore "gamecenter.mobi/paicode/client"
	"github.com/hyperledger/fabric/peerex"
)

const defPaicodeName string = "gamepaicore_v01"
const defRegion string = "gamepai01"

var defClient *clicore.ClientCore 

var mainCmd = &cobra.Command{
	Use: "gamepai [listeningaddr]",
	
	PreRunE: func(cmd *cobra.Command, args []string) error {

		config := &peerex.GlobalConfig{}
		//TODO: apply log to file
		err := config.InitGlobal(true)
		
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
			defClient.SetRpcRegion(defRegion)
			defClient.Rpc.Rpcbuilder.ChaincodeName = defPaicodeName
			restLogger.Infof("Start rpc, chaincode is %s", defClient.Rpc.Rpcbuilder.ChaincodeName)
				
		}else{
			restLogger.Info("Run under off-line mode")
		}
		
		return nil

	},
	
	Run: func(cmd *cobra.Command, args []string){
		
		if listenaddr == ""{
			listenaddr = "localhost:7280"
		}
		
		defClient.Accounts.KeyMgr.Load()
		//defer defClient.Accounts.KeyMgr.Persist()			
		
		// Initialize the REST service object
		restLogger.Infof("Initializing the REST service on %s", listenaddr)
	
		router := buildRouter()
		err := http.ListenAndServe(listenaddr, router)
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
var logtostd bool = false
var listenaddr string = ""

func main() {
	
	mainCmd.Flags().BoolVar(&debugmode, "debug", false, "run http server with debug output")
	mainCmd.Flags().BoolVar(&offlinemode, "offline", false, "not communicate with other peers")
	mainCmd.Flags().BoolVar(&logtostd, "logtostd", false, "put log to std out")
	mainCmd.Flags().StringVar(&listenaddr, "listen", "", "set listening addr")
	
	mainCmd.AddCommand(exitCmd)	

	err := mainCmd.Execute()
	if err != nil{
		fmt.Println("Command handler error:", err)
		os.Exit(1)		
	}

}


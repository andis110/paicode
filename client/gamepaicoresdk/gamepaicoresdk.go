package gamepaicoresdk // import "gamecenter.mobi/paicode/client/gamepaicoresdk"

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gocraft/web"
	"github.com/hyperledger/fabric/peerex"

	clicore "gamecenter.mobi/paicode/client"
	gamepaicorecommon "gamecenter.mobi/paicode/client/gamepaicorecommon"
)

const defPaicodeName string = "gamepaicore_v01"
const defRegion string = "gamepai01"

var restLogger = peerex.InitLogger("gamepaiREST")

//var debugmode bool = false
var offlinemode bool = false

//var logtostd bool = false
var listenaddr string = ""
var router *web.Router

type GamepaiCoreConfig struct {
	FileSystemPath string
	CrtFileName    string
	YamlFileName   string
	Address        string
	Port           int
}

func StartCoreDaemon(config string) string {
	log.Println("config: ", config)

	var coreConfig GamepaiCoreConfig
	err := json.Unmarshal([]byte(config), &coreConfig)
	if err != nil {
		log.Println("Parse config error: ", err)
		return fmt.Sprintf("failed.")
	}

	log.Println("FileSystemPath: ", coreConfig.FileSystemPath)
	log.Println("CrtFileName: ", coreConfig.CrtFileName)
	log.Println("YamlFileName: ", coreConfig.YamlFileName)
	log.Println("Address: ", coreConfig.Address)
	log.Println("Port: ", coreConfig.Port)
	crtFile := filepath.Join(coreConfig.FileSystemPath, coreConfig.CrtFileName+".crt")
	yamlFile := filepath.Join(coreConfig.FileSystemPath, coreConfig.YamlFileName+".yaml")
	log.Println("CrtFile: ", crtFile)
	log.Println("YamlFile: ", yamlFile)

	globalConfig := &peerex.GlobalConfig{}
	globalConfig.ConfigPath = make([]string, 1, 10)
	globalConfig.ConfigPath[0] = coreConfig.FileSystemPath // Path to look for the config file in
	globalConfig.ConfigFileName = coreConfig.YamlFileName

	defaultViperSetting := make(map[string]string)
	defaultViperSetting["peer.fileSystemPath"] = coreConfig.FileSystemPath
	defaultViperSetting["peer.tls.rootcert.file"] = crtFile

	err = globalConfig.InitGlobalWrapper(true, defaultViperSetting)
	if err != nil {
		log.Println("Init global config error: ", err)
		return fmt.Sprintf("failed.")
	}

	err = os.MkdirAll(globalConfig.GetPeerFS(), 0777)
	if err != nil {
		restLogger.Error("Mkdir error: ", err)
		return fmt.Sprintf("failed.")
	}

	gamepaicorecommon.DefClient = clicore.NewClientCore(globalConfig)

	if !offlinemode {

		conn := peerex.ClientConn{}

		err := conn.Dialdefault()
		if err != nil {
			restLogger.Error("Dial default error: ", err)
			return fmt.Sprintf("failed.")
		}

		gamepaicorecommon.DefClient.PrepareRpc(conn)
		gamepaicorecommon.DefClient.SetRpcRegion(defRegion)
		gamepaicorecommon.DefClient.Rpc.Rpcbuilder.ChaincodeName = defPaicodeName
		restLogger.Infof("Start rpc, chaincode is %s", gamepaicorecommon.DefClient.Rpc.Rpcbuilder.ChaincodeName)
	} else {
		restLogger.Info("Run under off-line mode")
	}

	if listenaddr == "" {
		listenaddr = fmt.Sprintf("%s:%d", coreConfig.Address, coreConfig.Port)
		//listenaddr = "0.0.0.0:7280"
	}

	gamepaicorecommon.DefClient.Accounts.KeyMgr.Load()
	//defer defClient.Accounts.KeyMgr.Persist()

	// Initialize the REST service object
	restLogger.Infof("Initializing the REST service on %s", listenaddr)
	router = gamepaicorecommon.BuildRouter()
	go startHttpServer()

	return "success"
}

func startHttpServer() {
	restLogger.Info("Starting HTTP Server ...")
	err := http.ListenAndServe(listenaddr, router)
	restLogger.Info("HTTP server is stopped.")
	if err != nil {
		restLogger.Error("Listen and Serve error: ", err)
	}

	if gamepaicorecommon.DefClient.IsRpcReady() {
		gamepaicorecommon.DefClient.ReleaseRpc()
	}
}

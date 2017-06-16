package main

import (
	"os"
	"fmt"
	_ "bufio"
	_ "strings"
	
	clicore "gamecenter.mobi/paicode/client"
	
	"github.com/spf13/cobra"
	"github.com/hyperledger/fabric/peerex"
)

var mainCmd = &cobra.Command{
	Use: "client",
}

var defClient *clicore.ClientCore 

func main() {
	
	config := peerex.GlobalConfig{}
	err := config.InitGlobal()
	
	if err != nil{
		panic(err)		
	}
	
	defClient = clicore.NewClientCore()
	
	var default_conn peerex.ClientConn
	err = default_conn.Dialdefault()
	if err != nil{
		fmt.Println("Dial to peer fail:", err)
		os.Exit(1)
	}	
	
	defer default_conn.C.Close()
	
	
	
}
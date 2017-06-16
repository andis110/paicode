package main

import (
	"os"
	"fmt"
	"bufio"
	"strings"
	
	"github.com/spf13/cobra"
	"github.com/hyperledger/fabric/peerex"
)

var mainCmd = &cobra.Command{
	Use: "client",
}

func main() {
	
	config := peerex.GlobalConfig{}
	err := config.InitGlobal()
	
	if err != nil{
		panic(err)		
	}
	
	var default_conn peerex.ClientConn
	err = default_conn.Dialdefault()
	if err != nil{
		fmt.Println("Dial to peer fail:", err)
		os.Exit(1)
	}	
	
	defer default_conn.C.Close()
	
	
	
}
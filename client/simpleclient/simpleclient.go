package main

import (
	"os"
	"fmt"
	"bufio"
	_ "strings"
	
	arghelper "github.com/mattn/go-shellwords"
	clicore "gamecenter.mobi/paicode/client"
	
	"github.com/spf13/cobra"
	"github.com/hyperledger/fabric/peerex"
)

var mainCmd = &cobra.Command{
	Use: ">",
}

var defClient *clicore.ClientCore 
var default_conn peerex.ClientConn

func main() {
	
	config := peerex.GlobalConfig{}
	err := config.InitGlobal()
	
	if err != nil{
		panic(err)		
	}
	
	defClient = clicore.NewClientCore()
	
	defer func(){
		if default_conn.C != nil{
			default_conn.C.Close()	
		}		
	}()
	
	fmt.Print("Starting .... ")
	mainCmd.AddCommand(rpcCmd)
	mainCmd.AddCommand(accountCmd)
	
	mainCmd.SetArgs([]string{"help"})
	err = mainCmd.Execute()
	if err != nil{
		fmt.Println("Command handler error:", err)
		os.Exit(1)		
	}	
	
	
	reader := bufio.NewReader(os.Stdin)
	parser := arghelper.NewParser()
	
	var ln string = "" 
	
	for {
		retbyte, notfinished, err := reader.ReadLine()
		if err != nil{
			break
		}
		
		ln += string(retbyte)
		if !notfinished {
			//handle read command line
			//fmt.Println("We get:", ln)
			args, err := parser.Parse(ln)
			if err != nil{
				fmt.Println("Input parse error:", err)		
			}else{
				mainCmd.SetArgs(args)
				err = mainCmd.Execute()
				if err != nil{
					fmt.Println("Command error:", err)		
				}				
			}
			
			ln = ""
			fmt.Println("Continue to next command:")
		}
	}
	
	fmt.Println("Exiting ...")	
	
}
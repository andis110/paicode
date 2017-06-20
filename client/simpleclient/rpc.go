package main

import (
	"fmt"
	_ "errors"
	
	"github.com/spf13/cobra"
	"github.com/hyperledger/fabric/peerex"
)

var rpcCmd = &cobra.Command{
	Use:   "rpc [command...]",
	Short: fmt.Sprintf("rpc commands."),
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error{
		
		if defClient.IsRpcReady(){
			return nil
		}
		
		conn := peerex.ClientConn{}
		err := conn.Dialdefault()
		if err != nil{
			return err
		}		
		
		defClient.PrepareRpc(conn) 
		return nil
	},
}

var userCmd = &cobra.Command{
	Use:   "user [command...]",
	Short: fmt.Sprintf("user commands."),
}

var registerCmd = &cobra.Command{
	Use:       "register",
	Short:     fmt.Sprintf("Register a public key"),
	RunE: func(cmd *cobra.Command, args []string) error{
		
		msg, err := defClient.Rpc.Registry(args...)
		if err != nil{
			return err
		}
		
		fmt.Print("Registry public key ok, TX id is", msg)
		return nil
	},
}


func init(){
	userCmd.AddCommand(registerCmd)
	
	rpcCmd.AddCommand(userCmd)
}
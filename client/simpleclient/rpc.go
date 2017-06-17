package main

import (
	"fmt"
	_ "errors"
	
	"github.com/spf13/cobra"
)

var rpcCmd = &cobra.Command{
	Use:   "rpc [command...]",
	Short: fmt.Sprintf("rpc commands."),
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error{
		
		if default_conn.C != nil{
			return nil
		}
		
		err := default_conn.Dialdefault()
		if err != nil{
			return err
		}
		return nil
	},
}

var userCmd = &cobra.Command{
	Use:   "user [command...]",
	Short: fmt.Sprintf("user commands."),
}

var registerCmd = &cobra.Command{
	Use:       "register <remark>",
	Short:     fmt.Sprintf("Register a public key"),
	RunE: func(cmd *cobra.Command, args []string) error{
		
		return nil
	},
}


func init(){
	userCmd.AddCommand(registerCmd)
	
	rpcCmd.AddCommand(userCmd)
}
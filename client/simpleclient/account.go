package main

import (
	"fmt"
	_ "errors"
	
	"github.com/spf13/cobra"
)

var accountCmd = &cobra.Command{
	Use:   "account [command...]",
	Short: fmt.Sprintf("account commands."),
}

var genPrivkeyCmd = &cobra.Command{
	Use:       "generate [remark]",
	Short:     fmt.Sprintf("generate a privkey"),
	Long:      fmt.Sprintf(`generate a privkey and save it with the name of remark.`),
	RunE: func(cmd *cobra.Command, args []string) error{
		
		return defClient.Accounts.GenPrivkey(args...)
	},
}

var dumpPrivkeyCmd = &cobra.Command{
	Use:       "dump [remark]",
	Short:     fmt.Sprintf("dump out a privkey"),
	RunE: func(cmd *cobra.Command, args []string) error{
		
		ret, err := defClient.Accounts.DumpPrivkey(args...)
		
		if err != nil{
			return err
		}
		
		fmt.Println(ret)		
		return nil	
	},
}

func init(){
	accountCmd.AddCommand(genPrivkeyCmd)
	accountCmd.AddCommand(dumpPrivkeyCmd)
}
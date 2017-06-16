package simpleclient

import (
	"fmt"
	"errors"
	
	"github.com/spf13/cobra"
	
	"gamecenter.mobi/paicode/client"
	"gamecenter.mobi/paicode/wallet"
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
		if len(args) > 1{
			return errors.New(fmt.Sprint("Could not recognize", args[1:]))
		}
		
		var remark string		
		if len(args) == 0{
			remark = client.RandStringRunes(16)
		}else{
			remark = args[0]
		}
		
		k, err := wallet.DefaultWallet.GeneratePrivKey()
		if err != nil{
			return err
		}
		
		defWallet.AddPrivKey(remark, k)		
	},
}

var dumpPrivkeyCmd = &cobra.Command{
	Use:       "dump [remark]",
	Short:     fmt.Sprintf("dump out a privkey"),
	RunE: func(cmd *cobra.Command, args []string) error{
		if len(args) != 1{
			return errors.New("Invalid remark")
		}
		
		k, err := defWallet.LoadPrivKey(args[0])
		if err != nil{
			return err
		}
		
		
		
		k, err := wallet.DefaultWallet.GeneratePrivKey()
		if err != nil{
			return err
		}
		
		defWallet.AddPrivKey(remark, k)		
	},
}

var defWallet *wallet.KeyManager = wallet.CreateSimpleManager("")

func init(){
	accountCmd.AddCommand(genPrivkeyCmd)
}
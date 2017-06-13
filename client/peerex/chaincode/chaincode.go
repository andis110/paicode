package chaincode

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/hyperledger/fabric/core"
	"github.com/hyperledger/fabric/peer/util"
	pb "github.com/hyperledger/fabric/protos"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	
	common "gamecenter.mobi/client/peerex" //temporary position
)

type chaincodeImpl struct{
	common.ClientConn
}

// attach a new client connection
func getDevopsClient(conn *common.ClientConn) (pb.DevopsClient) {

	devopsClient := pb.NewDevopsClient(conn.C)
	return devopsClient
}

// chaincodeInvokeOrQuery invokes or queries the chaincode. If successful, the
// INVOKE form prints the transaction ID on STDOUT, and the QUERY form prints
// the query result on STDOUT. A command-line flag (-r, --raw) determines
// whether the query result is output as raw bytes, or as a printable string.
// The printable form is optionally (-x, --hex) a hexadecimal representation
// of the query response. If the query response is NIL, nothing is output.
func (impl *chaincodeImpl) chaincodeInvokeOrQuery(cmd *cobra.Command, args []string, invoke bool) (err error) {
	spec, err := getChaincodeSpecification(cmd)
	if err != nil {
		return err
	}

	devopsClient := getDevopsClient(&impl.ClientConn)

	// Build the ChaincodeInvocationSpec message
	invocation := &pb.ChaincodeInvocationSpec{ChaincodeSpec: spec}
	if customIDGenAlg != common.UndefinedParamValue {
		invocation.IdGenerationAlg = customIDGenAlg
	}

	var resp *pb.Response
	if invoke {
		resp, err = devopsClient.Invoke(context.Background(), invocation)
	} else {
		resp, err = devopsClient.Query(context.Background(), invocation)
	}

	if err != nil {
		if invoke {
			err = fmt.Errorf("Error invoking %s: %s\n", chainFuncName, err)
		} else {
			err = fmt.Errorf("Error querying %s: %s\n", chainFuncName, err)
		}
		return
	}
	if invoke {
		transactionID := string(resp.Msg)
		logger.Infof("Successfully invoked transaction: %s(%s)", invocation, transactionID)
	} else {
		logger.Infof("Successfully queried transaction: %s", invocation)
		if resp != nil {
			if chaincodeQueryRaw {
				if chaincodeQueryHex {
					err = errors.New("Options --raw (-r) and --hex (-x) are not compatible\n")
					return
				}
				fmt.Print("Query Result (Raw): ")
				os.Stdout.Write(resp.Msg)
			} else {
				if chaincodeQueryHex {
					fmt.Printf("Query Result: %x\n", resp.Msg)
				} else {
					fmt.Printf("Query Result: %s\n", string(resp.Msg))
				}
			}
		}
	}
	return nil
}

func checkChaincodeCmdParams(cmd *cobra.Command) error {

	if chaincodeName == common.UndefinedParamValue {
		if chaincodePath == common.UndefinedParamValue {
			return fmt.Errorf("Must supply value for %s path parameter.\n", chainFuncName)
		}
	}

	// Check that non-empty chaincode parameters contain only Args as a key.
	// Type checking is done later when the JSON is actually unmarshaled
	// into a pb.ChaincodeInput. To better understand what's going
	// on here with JSON parsing see http://blog.golang.org/json-and-go -
	// Generic JSON with interface{}
	if chaincodeCtorJSON != "{}" {
		var f interface{}
		err := json.Unmarshal([]byte(chaincodeCtorJSON), &f)
		if err != nil {
			return fmt.Errorf("Chaincode argument error: %s", err)
		}
		m := f.(map[string]interface{})
		sm := make(map[string]interface{})
		for k := range m {
			sm[strings.ToLower(k)] = m[k]
		}
		_, argsPresent := sm["args"]
		_, funcPresent := sm["function"]
		if !argsPresent || (len(m) == 2 && !funcPresent) || len(m) > 2 {
			return fmt.Errorf("Non-empty JSON chaincode parameters must contain the following keys: 'Args' or 'Function' and 'Args'")
		}
	} else {
		return errors.New("Empty JSON chaincode parameters must contain the following keys: 'Args' or 'Function' and 'Args'")
	}

	if chaincodeAttributesJSON != "[]" {
		var f interface{}
		err := json.Unmarshal([]byte(chaincodeAttributesJSON), &f)
		if err != nil {
			return fmt.Errorf("Chaincode argument error: %s", err)
		}
	}

	return nil
}


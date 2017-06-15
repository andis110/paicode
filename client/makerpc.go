package client

import (
	fabric_pb "github.com/hyperledger/fabric/protos"
	"github.com/hyperledger/fabric/peerex"
)

type RpcBuilder struct{
	ChaincodeName	 string
//	ChaincodeLang    string
	Function		 string
	
	Security		 *SecurityPolicy
	Responding		 *RespondSepc
	
	Conn			 *peerex.ClientConn		
}

type SecurityPolicy struct{
	User			string
	Attributes		[]string
	Metadata		[]byte
	CustomIDGenAlg  string
}

type RespondSepc struct{
	QueryDumpHex	bool
}

var defaultSecPolicy = &SecurityPolicy{Attributes: []string{}}
var defaultRespond = &RespondSepc{}

func makeStringArgsToPb(funcname string, args []string) *fabric_pb.ChaincodeInput{
	
	input := &fabric_pb.ChaincodeInput{}
	//please remember the trick fabric used:
	//it push the "function name" as the first argument
	//in a rpc call
	var inarg [][]byte
	if len(funcname) == 0{
		input.Args = make([][]byte, len(args))	
		inarg = input.Args[:]
	}else{
		input.Args = make([][]byte, len(args) + 1)
		input.Args[0] = []byte(funcname)
		inarg = input.Args[1:]
	}
	
	for i, arg := range args{
		inarg[i] = []byte(arg)
	}
	
	return input
}

func (b *RpcBuilder) Fire(id string, args []string) error{
	
	spec := &fabric_pb.ChaincodeSpec{
		Type: fabric_pb.ChaincodeSpec_GOLANG,	//always set it as golang
		ChaincodeID: &fabric_pb.ChaincodeID{Name: b.ChaincodeName},
		CtorMsg : makeStringArgsToPb(b.Function, args),
	}
	
	invocation := &fabric_pb.ChaincodeInvocationSpec{ChaincodeSpec: spec}
	
	if b.Security != nil{
		spec.Attributes = b.Security.Attributes
		if len(b.Security.CustomIDGenAlg) != 0{
			invocation.IdGenerationAlg = b.Security.CustomIDGenAlg
		}
	}
	
	//final check attributes
	if spec.Attributes == nil{
		spec.Attributes = defaultSecPolicy.Attributes
	}	
	
	resp, err := fabric_pb.DevopsClient(b.Conn).Invoke(ctx, invocation)
	
	if err != nil{
		return err
	}
	
	return nil
}

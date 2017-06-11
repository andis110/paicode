package security

import(
	
	"strings"
	"github.com/hyperledger/fabric/core/chaincode/shim"	
	
	pb "gamecenter.mobi/paicode/protos" 
)

type SecurityPolicy struct{
	dbgMode bool
	netCode int
}

const(
		
	AdminPrivilege string = "Admin"
	ManagerPrivilege string = "Manager"
	DelegatePrivilege string = "Delegate"
	ObserverPrivilege string = "Observer"
	
	debugPrivilege string = "debug"
	noPrivilege string = "none"
	
	privilege_Attr string = "PaiAdminRole"
	region_Attr string = "PaiAdminRegion"
	
	debugRegion string = "debug"
	noRegion string = "none"
)

//keep a singleton
var Helper *SecurityPolicy

func InitSecHelper(set *pb.DeploySetting) *SecurityPolicy{
	
	Helper = &SecurityPolicy{dbgMode: set.DebugMode, netCode: int(set.NetworkCode)}
	return Helper
}

func (sec *SecurityPolicy) ActiveAudit(stub shim.ChaincodeStubInterface, desc string){
	
}

func (sec *SecurityPolicy) GetPrivilege(stub shim.ChaincodeStubInterface) (privilege string, region string){
	
	if sec.dbgMode{
		return debugPrivilege, debugRegion
	}
	
	cert, err := stub.GetCallerCertificate()
	if err != nil || cert == nil{
		return noPrivilege, noRegion
	}
	
	
	return noPrivilege, noRegion
	
}

func (sec *SecurityPolicy) VerifyPrivilege(certpriv string, expect string) bool{
	
	if sec.dbgMode && strings.Compare(certpriv, debugPrivilege) == 0{
		return true
	}
		
	return strings.Compare(certpriv, expect) != 0
	
}

func (sec *SecurityPolicy) VerifyRegion(region string, expect string) bool{
	
	if sec.dbgMode && strings.Compare(region, debugRegion) == 0{
		return true
	}	
	
	return strings.Compare(region, expect) != 0
	
}

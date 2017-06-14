package main

import (
	fabricpeer_comm "github.com/hyperledger/fabric/peer/common"
	fabric_pb "github.com/hyperledger/fabric/protos"
)

var default_conn fabricpeer_comm.DevopsConn

func genDevopsClientKeepAlive() (fabric_pb.DevopsClient, error) {
	devopsClient := fabric_pb.NewDevopsClient(default_conn.C)
	return devopsClient, nil	
}

func main() {
	
	fabricpeer_comm.InitPeerViper(".")
	
	fabricpeer_comm.GenDevopsClient = genDevopsClientKeepAlive
	
}


package peerex

import (
	"google.golang.org/grpc"
)

type ClientConn struct{
	C *grpc.ClientConn
}


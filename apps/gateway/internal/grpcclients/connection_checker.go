package grpcclients

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ConnectionChecker(cc *grpc.ClientConn) error {
	if cc == nil {
		return status.Error(codes.Unavailable, "nil grpc ClientConn (downstream not connected)")
	}
	return nil
}

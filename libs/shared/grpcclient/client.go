package grpcclient

import (
	"fmt"

	"google.golang.org/grpc"
)

type Client struct {
	connections map[string]*grpc.ClientConn
	dialOps     []grpc.DialOption
}

func NewClient(dialOps ...grpc.DialOption) *Client {
	return &Client{
		connections: make(map[string]*grpc.ClientConn),
		dialOps:     dialOps,
	}
}

func (c *Client) GetConnection(addr string) (*grpc.ClientConn, error) {
	if conn, ok := c.connections[addr]; ok {
		return conn, nil
	}
	conn, err := grpc.Dial(addr, c.dialOps...)
	if err != nil {
		return nil, fmt.Errorf("Failed to dial gRPC Server %s: %w", addr, err)
	}
	c.connections[addr] = conn
	return conn, nil
}

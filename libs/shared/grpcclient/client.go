// Package grpcclient provides reusable gRPC client logic for connecting to services.
package grpcclient

import (
	"fmt"

	"google.golang.org/grpc"
)

// Client is a reusable gRPC client that manages connections to gRPC servers.
type Client struct {
	connections map[string]*grpc.ClientConn
	dialOps     []grpc.DialOption
}

// NewClient returns a new internal gRPC client with optional dial options.
func NewClient(dialOps ...grpc.DialOption) *Client {
	return &Client{
		connections: make(map[string]*grpc.ClientConn),
		dialOps:     dialOps,
	}
}

// GetConnection returns a cached or newly established gRPC connection for the given address.
func (c *Client) GetConnection(addr string) (*grpc.ClientConn, error) {
	if conn, ok := c.connections[addr]; ok {
		return conn, nil
	}

	conn, err := grpc.NewClient(addr, c.dialOps...)
	if err != nil {
		return nil, fmt.Errorf("failed to dial gRPC Server %s: %w", addr, err)
	}
	c.connections[addr] = conn
	return conn, nil
}

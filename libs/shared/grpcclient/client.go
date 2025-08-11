package grpcclient

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ToxicToast/Azkaban-Go/libs/shared/helper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type Client struct {
	mu    sync.Mutex
	conns map[string]*grpc.ClientConn // key = service name
	cfg   map[string]helper.ServiceCfg
	idx   map[string]int
}

func NewClient(services map[string]helper.ServiceCfg) *Client {
	return &Client{
		conns: map[string]*grpc.ClientConn{},
		cfg:   services,
		idx:   make(map[string]int),
	}
}

func (c *Client) Get(ctx context.Context, service string) (*grpc.ClientConn, error) {
	// 1) Cache-Hit?
	c.mu.Lock()
	if cc, ok := c.conns[service]; ok && cc != nil {
		c.mu.Unlock()
		return cc, nil
	}

	// 2) Config + Endpoint wählen (Round-Robin)
	sc, ok := c.cfg[service]
	if !ok || len(sc.Endpoints) == 0 {
		c.mu.Unlock()
		return nil, fmt.Errorf("unknown service %q or no endpoints", service)
	}
	i := c.idx[service] % len(sc.Endpoints)
	c.idx[service]++
	target := sc.Endpoints[i]
	c.mu.Unlock()

	// 3) Dial (blocking, mit Timeout übers ctx)
	var opts []grpc.DialOption
	if sc.Insecure {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	// TODO: TLS creds bauen (CA, ServerName, mTLS etc.)

	cc, err := grpc.NewClient(target, opts...)
	if err != nil {
		return nil, fmt.Errorf("dial %s: %w", target, err)
	}

	// 4) Conn speichern (kurz locken) und zurück
	c.mu.Lock()
	if _, exists := c.conns[service]; !exists {
		c.conns[service] = cc
	}
	c.mu.Unlock()

	return cc, nil
}

func (c *Client) Ping(ctx context.Context, service string) error {
	cc, err := c.Get(ctx, service)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	_, err = healthpb.NewHealthClient(cc).Check(ctx, &healthpb.HealthCheckRequest{Service: service})
	return err
}

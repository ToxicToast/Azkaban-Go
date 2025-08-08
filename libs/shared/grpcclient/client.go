package grpcclient

import (
	"context"
	"fmt"
	"sync"

	"github.com/ToxicToast/Azkaban-Go/libs/shared/helper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

/*func (c *Client) Get(service string) (*grpc.ClientConn, error) {
	c.mu.Lock()
	if cc, ok := c.conns[service]; ok {
		c.mu.Unlock()
		log.Printf("downstream cache hit service=%s conn=%p state=%s", service, cc, cc.GetState().String())
		return cc, nil
	}

	log.Printf("dialing service=%s; connection=%p", service, c.conns[service])

	sc, ok := c.cfg[service]
	if !ok || len(sc.Endpoints) == 0 {
		c.mu.Unlock()
		return nil, fmt.Errorf("unknown service %q or no endpoints (%s)", service, len(sc.Endpoints))
	}

	target := sc.Endpoints[0]

	log.Printf("dialing downstream %q for service=%s", target, service)

	opts := []grpc.DialOption{
		grpc.WithBlock(),
	}
	if sc.Insecure {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		// ca_file, server_name etc.
	}

	to := time.Duration(sc.DialTimeoutms) * time.Millisecond
	if to <= 0 {
		to = 3 * time.Second
	}
	cc, err := grpc.NewClient(target, opts...)
	if err != nil {
		return nil, fmt.Errorf("dial %s: %w", target, err)
	}

	log.Printf("downstream connected service=%s target=%s conn=%p state=%s",
		service, target, cc, cc.GetState().String())

	c.mu.Lock()
	c.conns[service] = cc
	c.mu.Unlock()
	return cc, nil
}*/

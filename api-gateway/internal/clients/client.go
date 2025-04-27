package clients

import (
	"context"
	"fmt"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClientPool struct {
	conn     *grpc.ClientConn
	lastUsed time.Time
	mutex    sync.Mutex
	addr     string
	maxIdle  time.Duration
}

func NewGrpcClientPool(addr string) *GrpcClientPool {
	return &GrpcClientPool{
		addr:    addr,
		maxIdle: 5 * time.Minute,
	}
}

func (p *GrpcClientPool) GetConn() (*grpc.ClientConn, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.conn != nil && p.conn.GetState() == connectivity.Ready {
		p.lastUsed = time.Now()
		return p.conn, nil
	}

	if p.conn != nil {
		p.conn.Close()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, p.addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %w", p.addr, err)
	}

	p.conn = conn
	p.lastUsed = time.Now()
	return conn, nil
}

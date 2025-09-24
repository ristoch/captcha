package grpc

import (
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ConnectionPool struct {
	address       string
	pool          []*grpc.ClientConn
	mu            sync.RWMutex
	maxSize       int
	current       int
	cleanupTicker *time.Ticker
	stopChan      chan struct{}
}

func NewConnectionPool(address string, maxSize int) *ConnectionPool {
	pool := &ConnectionPool{
		address:  address,
		pool:     make([]*grpc.ClientConn, 0, maxSize),
		maxSize:  maxSize,
		stopChan: make(chan struct{}),
	}

	pool.startCleanup()
	return pool
}

func (p *ConnectionPool) GetConnection() (*grpc.ClientConn, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.current < len(p.pool) && p.pool[p.current] != nil {
		conn := p.pool[p.current]
		p.current++
		return conn, nil
	}

	if len(p.pool) < p.maxSize {
		conn, err := grpc.Dial(p.address, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, err
		}

		p.pool = append(p.pool, conn)
		p.current++
		return conn, nil
	}

	conn := p.pool[0]
	p.current = 1
	return conn, nil
}

func (p *ConnectionPool) ReturnConnection(conn *grpc.ClientConn) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.current > 0 {
		p.current--
	}
}

func (p *ConnectionPool) startCleanup() {
	p.cleanupTicker = time.NewTicker(10 * time.Minute)

	go func() {
		for {
			select {
			case <-p.cleanupTicker.C:
				p.cleanup()
			case <-p.stopChan:
				return
			}
		}
	}()
}

func (p *ConnectionPool) cleanup() {
	p.mu.Lock()
	defer p.mu.Unlock()

	keepCount := len(p.pool) / 2
	if keepCount < 1 {
		keepCount = 1
	}

	for i := keepCount; i < len(p.pool); i++ {
		if p.pool[i] != nil {
			p.pool[i].Close()
		}
	}

	p.pool = p.pool[:keepCount]
	p.current = 0
}

func (p *ConnectionPool) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.cleanupTicker != nil {
		p.cleanupTicker.Stop()
	}

	for _, conn := range p.pool {
		if conn != nil {
			conn.Close()
		}
	}

	close(p.stopChan)
}

func (p *ConnectionPool) GetStats() map[string]interface{} {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return map[string]interface{}{
		"total_connections": len(p.pool),
		"max_size":          p.maxSize,
		"current_index":     p.current,
		"usage_percent":     float64(p.current) / float64(p.maxSize) * 100,
	}
}

package grpc

import (
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ConnectionPool manages a pool of gRPC connections
type ConnectionPool struct {
	address       string
	pool          []*grpc.ClientConn
	mu            sync.RWMutex
	maxSize       int
	current       int
	cleanupTicker *time.Ticker
	stopChan      chan struct{}
}

// NewConnectionPool creates a new connection pool
func NewConnectionPool(address string, maxSize int) *ConnectionPool {
	pool := &ConnectionPool{
		address:  address,
		pool:     make([]*grpc.ClientConn, 0, maxSize),
		maxSize:  maxSize,
		stopChan: make(chan struct{}),
	}

	// Start cleanup goroutine
	pool.startCleanup()
	return pool
}

// GetConnection returns a connection from the pool
func (p *ConnectionPool) GetConnection() (*grpc.ClientConn, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Try to reuse existing connection
	if p.current < len(p.pool) && p.pool[p.current] != nil {
		conn := p.pool[p.current]
		p.current++
		return conn, nil
	}

	// Create new connection if pool is not full
	if len(p.pool) < p.maxSize {
		conn, err := grpc.Dial(p.address, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, err
		}

		p.pool = append(p.pool, conn)
		p.current++
		return conn, nil
	}

	// Pool is full, reuse oldest connection
	conn := p.pool[0]
	p.current = 1
	return conn, nil
}

// ReturnConnection returns a connection to the pool
func (p *ConnectionPool) ReturnConnection(conn *grpc.ClientConn) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Reset current pointer for reuse
	if p.current > 0 {
		p.current--
	}
}

// startCleanup starts a background goroutine to clean up idle connections
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

// cleanup removes idle connections
func (p *ConnectionPool) cleanup() {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Keep only half of the connections to reduce memory usage
	keepCount := len(p.pool) / 2
	if keepCount < 1 {
		keepCount = 1
	}

	// Close excess connections
	for i := keepCount; i < len(p.pool); i++ {
		if p.pool[i] != nil {
			p.pool[i].Close()
		}
	}

	// Resize pool
	p.pool = p.pool[:keepCount]
	p.current = 0
}

// Close closes all connections in the pool
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

// GetStats returns pool statistics
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

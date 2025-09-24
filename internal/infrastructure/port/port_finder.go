package port

import (
	"fmt"
	"net"
	"time"
)

type PortFinder struct {
	minPort int
	maxPort int
}

func NewPortFinder(minPort, maxPort int) *PortFinder {
	return &PortFinder{
		minPort: minPort,
		maxPort: maxPort,
	}
}

func (p *PortFinder) FindAvailablePortWithRetry(maxRetries int, retryInterval time.Duration) (int, error) {
	for i := 0; i < maxRetries; i++ {
		port, err := p.FindAvailablePort()
		if err == nil {
			return port, nil
		}

		if i < maxRetries-1 {
			time.Sleep(retryInterval)
		}
	}

	return 0, fmt.Errorf("failed to find available port after %d retries", maxRetries)
}

func (p *PortFinder) FindAvailablePort() (int, error) {
	for port := p.minPort; port <= p.maxPort; port++ {
		addr := fmt.Sprintf(":%d", port)
		listener, err := net.Listen("tcp", addr)
		if err == nil {
			listener.Close()
			return port, nil
		}
	}
	return 0, fmt.Errorf("no available port found in range %d-%d", p.minPort, p.maxPort)
}

package utils

import (
	"fmt"
	"net"
	"time"
)

func FindAvailablePort(minPort, maxPort int32) (int32, error) {
	for port := minPort; port <= maxPort; port++ {
		address := fmt.Sprintf(":%d", port)
		listener, err := net.Listen("tcp", address)
		if err == nil {
			listener.Close()
			return port, nil
		}
	}
	return 0, fmt.Errorf("no available port found in range %d-%d", minPort, maxPort)
}

func IsPortAvailable(port int32) bool {
	address := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return false
	}
	listener.Close()
	return true
}

func WaitForPort(host string, port int32, timeout time.Duration) error {
	address := fmt.Sprintf("%s:%d", host, port)
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", address, time.Second)
		if err == nil {
			conn.Close()
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}

	return fmt.Errorf("port %d on %s not available after %v", port, host, timeout)
}

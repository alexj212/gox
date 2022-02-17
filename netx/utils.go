package netx

import (
    "fmt"
    "github.com/pkg/errors"
    "net"
)

// GetFreePort asks the kernel for a free open port that is ready to use.
func GetFreePort(start, end int) (int, error) {
    if start < 1 || start > 65535 {
        return -1, errors.Errorf("invalid start port: %d", start)
    }

    if end < 1 || end > 65535 {
        return -1, errors.Errorf("invalid end port: %d", end)
    }
    if start >= end {
        return -1, errors.Errorf("invalid start: %d/end: %d  portS", start, end)
    }

    var port = -1
    for p := start; p < end; p++ {
        avail, err := IsPortFree(p)
        if err == nil && avail == true {
            port = p
            break
        }
    }
    if port == -1 {
        return -1, errors.Errorf("unable to find port in range[%d:%d]", start, end)
    }
    return port, nil
}

// IsPortFree asks the kernel is a specific port is free
func IsPortFree(port int) (bool, error) {
    if port < 1 || port > 65535 {
        return false, errors.Errorf("invalid port: %d", port)
    }
    listener := fmt.Sprintf(":%d", port)
    addr, err := net.ResolveTCPAddr("tcp", listener)
    if err != nil {
        return false, err
    }

    l, err := net.ListenTCP("tcp", addr)
    if err != nil {
        return false, err
    }
    defer func() {
        _ = l.Close()
    }()
    return true, nil
}

// GetRandomFreePort asks the kernel for a free open port that is ready to use.
func GetRandomFreePort() (int, error) {
    addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
    if err != nil {
        return 0, err
    }

    l, err := net.ListenTCP("tcp", addr)
    if err != nil {
        return 0, err
    }
    defer func() {
        _ = l.Close()
    }()
    return l.Addr().(*net.TCPAddr).Port, nil
}

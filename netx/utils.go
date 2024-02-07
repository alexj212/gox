package netx

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/pkg/errors"
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

func fetchIp(url string) (string, error) {
	client := http.Client{
		Timeout: time.Second * 1,
	}

	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		return string(bodyBytes), nil
	}
	return "", fmt.Errorf("error fetching ip: %s", resp.Status)
}

// GetLocalIP returns the non loopback local IP of the host
func GetLocalIP() string {
	return GetIP()
}

// GetOutboundIP Get preferred outbound ip of this machine
func GetOutboundIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.To4().String(), nil
}

// GetManagementIp return env value for MANAGEMENT_IP
func GetManagementIp() string {
	ip := os.Getenv("MANAGEMENT_IP")

	if ip == "" {
		ip = os.Getenv("CLIENT_IP")
	}
	if ip == "" {
		ip = os.Getenv("HOST_IP")
	}
	if ip == "" {
		ip = GetIP()
	}

	return ip
}

// GetClientIp return env value for CLIENT_IP
func GetClientIp() string {
	ip := os.Getenv("CLIENT_IP")

	if ip == "" {
		ip = os.Getenv("MANAGEMENT_IP")
	}
	if ip == "" {
		ip = os.Getenv("HOST_IP")
	}
	if ip == "" {
		ip = GetIP()
	}

	return ip
}

// GetHostIp return env value for HOST_IP
func GetHostIp() string {
	ip := os.Getenv("HOST_IP")

	if ip == "" {
		ip = os.Getenv("MANAGEMENT_IP")
	}
	if ip == "" {
		ip = os.Getenv("CLIENT_IP")
	}
	if ip == "" {
		ip = GetIP()
	}

	return ip
}

func GetIP() string {
	resp, err := fetchIp("http://169.254.169.254/latest/meta-data/local-ipv4")
	if err == nil {
		return resp
	}

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

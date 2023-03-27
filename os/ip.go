package os

import "net"

var localIP string

// GetLocalIP 获取本地ip
func GetLocalIP() string {
	if localIP != "" {
		return localIP
	}
	localhost := "127.0.0.1"
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		localIP = localhost
		return localIP
	}
	for _, address := range addresses {
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				localIP = ipNet.IP.String()
				return localIP
			}
		}
	}
	localIP = localhost
	return localIP
}

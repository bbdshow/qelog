package kit

import (
	"net"
	"strings"
)

func GetLocalIPV4() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", nil
	}
	for _, addr := range addrs {
		// 检查ip地址是否为回环地址
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", nil
}

func GetLocalIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		return "", err
	}
	if updAddr, ok := conn.LocalAddr().(*net.UDPAddr); ok {
		if updAddr.IP.To4() != nil || updAddr.IP.To16() != nil {
			return updAddr.IP.String(), nil
		}
	}
	return AddrStringToIP(conn.LocalAddr()), nil
}

func AddrStringToIP(addr net.Addr) string {
	str := addr.String()
	if strings.HasPrefix(str, "[") {
		// ipv6
		return strings.Split(strings.TrimPrefix(str, "["), "]:")[0]
	}
	return strings.Split(str, ":")[0]
}

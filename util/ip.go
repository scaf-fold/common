package util

import (
	"fmt"
	"net"
	"strings"
)

func Ip(addr chan string) {
	addrs, err := net.InterfaceAddrs()
	defer close(addr)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				addr <- ipnet.IP.String()
				return
			}
		}
	}
	addr <- "127.0.0.1"
}

func PublicIp(addr chan string) {
	conn, err := net.Dial("udp", "8.8.8.8:53")
	defer close(addr)
	if err != nil {
		fmt.Println(err)
		addr <- "127.0.0.1"
		return
	}
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	ip := strings.Split(localAddr.String(), ":")[0]
	addr <- ip
}

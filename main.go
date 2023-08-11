package main

import (
	"fmt"
	"net"
	"time"
)

type Device struct {
	IP    net.IP
	Ports []int
}

func scanIPRange(startIP, endIP net.IP, timeout time.Duration, ports []int) {
	fmt.Printf("Scanning IP range %s - %s\n", startIP, endIP)

	for ip := startIP; ip4Cmp(ip, endIP) != 0; incrementIP(ip) {
		conn, err := net.DialTimeout("ip4:icmp", ip.String(), timeout)
		if err == nil {
			// fmt.Printf("Device found: %s\n", ip)
			scanPorts(Device{IP: ip, Ports: ports}, timeout)
			conn.Close()
		}
	}

	// Добавляем последний IP-адрес endIP в список устройств
	scanPorts(Device{IP: endIP, Ports: ports}, timeout)
}

func ip4Cmp(ip1, ip2 net.IP) int {
	ip1 = ip1.To4()
	ip2 = ip2.To4()
	for i := 0; i < 4; i++ {
		if ip1[i] < ip2[i] {
			return -1
		} else if ip1[i] > ip2[i] {
			return 1
		}
	}
	return 0
}

func incrementIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func scanPorts(device Device, timeout time.Duration) {
	for _, port := range device.Ports {
		target := fmt.Sprintf("%s:%d", device.IP.String(), port)
		conn, err := net.DialTimeout("tcp", target, timeout)
		if err == nil {
			fmt.Printf("Port %d is open on device %s\n", port, device.IP)
			conn.Close()
		}
	}
}

func main() {
	startIP := net.ParseIP("91.105.192.100")
	endIP := net.ParseIP("91.105.192.101")
	timeout := 1 * time.Second
	ports := []int{80, 443, 8080, 22, 2000, 5000}

	scanIPRange(startIP, endIP, timeout, ports)

	fmt.Println("Scanning complete.")
}

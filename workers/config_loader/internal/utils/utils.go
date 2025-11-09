package utils

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

// ExpandCIDR returns all IP addresses in a CIDR range without network and broadcast addresses
// Example: ExpandCIDR("38.0.101.0/30") returns []string{"38.0.101.1", "38.0.101.2"}
func ExpandCIDR(cidr string) ([]string, error) {
	_, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, fmt.Errorf("invalid cidr")
	}

	var ips []string
	for ip := nextIP(ipnet.IP.Mask(ipnet.Mask)); ipnet.Contains(ip); ip = nextIP(ip) {
		ips = append(ips, ip.String())
	}
	return ips[:len(ips)-1], nil
}

func nextIP(ip net.IP) net.IP {
	ip = append(net.IP(nil), ip...)
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] != 0 {
			break
		}
	}
	return ip
}

// ExpandRange returns all IP addresses in range
// Example: ExpandRange("38.0.100.255-38.101.1") returns []string{"38.0.100.255", "38.0.101.0", "38.101.1"}
func ExpandRange(iprange string) ([]string, error) {
	parts := strings.SplitN(iprange, "-", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid ip range")
	}

	startIP, endIP := parts[0], parts[1]
	s := net.ParseIP(startIP).To4()
	if s == nil {
		return nil, fmt.Errorf("invalid start range IP")
	}
	e := net.ParseIP(endIP).To4()
	if e == nil {
		return nil, fmt.Errorf("invalid end range IP")
	}

	var ips []string
	for ip := append(net.IP(nil), s...); !ipAfter(ip, e); ip = nextIP(ip) {
		ips = append(ips, ip.String())
	}
	return ips, nil
}

func ipAfter(a, b net.IP) bool {
	for i := 0; i < 4; i++ {
		if a[i] > b[i] {
			return true
		}
		if a[i] < b[i] {
			return false
		}
	}
	return false
}

// ExpandIpFromN returns all IP addresses by template with variables {X}, {Y}
// Example: ExpandIpFromN(0, 6, 32, 1, 3, "10.{X}.{Y}.5") returns []string{"10.32.1.5", "10.32.2.5", "10.33.1.5", "10.33.2.5", "10.34.1.5"}
func ExpandIpFromN(nStart, nEnd, offsetX, offsetY, block int, ipTmpl string) []string {
	ips := make([]string, 0, nEnd-nStart)
	for n := nStart; n < nEnd; n++ {
		ip := string([]byte(ipTmpl))
		ip = strings.Replace(
			ip, "{X}", strconv.Itoa(n/block+offsetX), 1,
		)
		ip = strings.Replace(
			ip, "{Y}", strconv.Itoa(n%block+offsetY), 1,
		)
		ips = append(ips, ip)
	}
	return ips
}

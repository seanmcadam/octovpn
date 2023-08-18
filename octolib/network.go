package octolib

import (
	"net"
	"regexp"
)

var ipv4Pattern *regexp.Regexp
var ipv6Pattern *regexp.Regexp

func init() {
	ipv4Pattern = regexp.MustCompile(`^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`)
	ipv6Pattern = regexp.MustCompile(`^([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}$`)
}

func ValidIP(ip string) bool {
	return (ValidIP4(ip) || ValidIP6(ip))
}

func ValidIP4(ip string) bool {
	return ipv4Pattern.MatchString(ip)
}

func ValidIP6(ip string) bool {
	return ipv6Pattern.MatchString(ip)
}

//
func ValidIPHost(host string) bool {
	ips, err := net.LookupIP(host)
	if err != nil {
		return false
	}

	// Check if at least one IP address is resolved
	return len(ips) > 0
}

//
func ValidIP4Host(host string) bool {
	ips, err := net.LookupIP(host)
	if err != nil {
		return false
	}

	// Check if at least one IP address is resolved
	if len(ips) > 0 {
	for _, address := range ips {
		if isValidIPv4Address(address.String()) {
			return true
		}
	}
	}
return false
}

//
func ValidIP6Host(host string) bool {
	ips, err := net.LookupIP(host)
	if err != nil {
		return false
	}

	// Check if at least one IP address is resolved
	if len(ips) > 0 {
	for _, address := range ips {
		if isValidIPv6Address(address.String()) {
			return true
		}
	}
	}

	return false
}


//
func isValidIPv4Address(address string) bool {
	ip := net.ParseIP(address)
	return ip != nil && ip.To4() != nil
}

//
func isValidIPv6Address(address string) bool {
	ip := net.ParseIP(address)
	return ip != nil && ip.To16() != nil
}

func isValidIPv6Host(host string) bool {
	addresses, err := net.LookupHost(host)
	if err != nil {
		return false
	}

	for _, address := range addresses {
		if isValidIPv6Address(address) {
			return true
		}
	}

	return false
}

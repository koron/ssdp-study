package udp

import (
	"errors"
	"fmt"
	"net"
	"strconv"
)

// Interface lookups net.Interface by query string.  Query string should be one
// of name, index (numeric string) or IP address.
func Interface(query string) (*net.Interface, error) {
	if ifi, err := net.InterfaceByName(query); err == nil {
		return ifi, nil
	}
	n, err := strconv.Atoi(query)
	if err == nil && n >= 0 {
		if ifi, err := net.InterfaceByIndex(n); err == nil {
			return ifi, nil
		}
	}
	ip := net.ParseIP(query)
	if ip != nil {
		if ifi, err := InterfacesByIP(ip); err == nil {
			return ifi, nil
		}
	}
	return nil, fmt.Errorf("no interfaces for %q", query)
}

// InterfacesByIP lookups net.Interfaces by IP address.
func InterfacesByIP(ip net.IP) (*net.Interface, error) {
	if ip == nil {
		return nil, errors.New("invalid nil IP")
	}
	list, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, ifi := range list {
		if contains(&ifi, ip) {
			return &ifi, nil
		}
	}
	return nil, fmt.Errorf("no interfaces for IP:%s", ip.String())
}

func contains(ifi *net.Interface, ip net.IP) bool {
	list, err := ifi.Addrs()
	if err != nil {
		return false
	}
	for _, addr := range list {
		if ipnet, ok := addr.(*net.IPNet); ok && ipnet.Contains(ip) {
			return true
		}
	}
	return false
}

// HasRealAddress checks net.Interface having real IP address or not.
func HasRealAddress(ifi *net.Interface) bool {
	addrs, err := ifi.Addrs()
	if err != nil {
		return false
	}
	for _, a := range addrs {
		ip := net.ParseIP(a.String())
		if !ip.IsUnspecified() {
			return true
		}
	}
	return false
}

package main

import (
	"fmt"
	"log"
	"net"
)

const (
	addrIP4 = "239.255.255.250:1900"
	addrIP6 = "[FF02::C]:1900"
)

const (
	stAll         = "ssdp:all"
	stRootDevices = "upnp:rootdevices"
)

func buildMSearch(st string, mx int) string {
	return fmt.Sprintf("M-SEARCH * HTTP/1.1\r\nHOST: %s\r\nMAN: \"ssdp:discover\"\r\nMX: %d\r\nST: %s\r\n\r\n", addrIP4, mx, st)
}

func main() {
	addr, err := net.ResolveUDPAddr("udp", addrIP4)
	if err != nil {
		log.Fatalf("net.ResolveUDPAddr() failed: %s", err)
	}

	connReq, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatalf("net.DialUDP() failed: %s", err)
	}
	sec := 10
	defer connReq.Close()

	// send
	msg := buildMSearch(stAll, sec)
	if _, err := connReq.Write([]byte(msg)); err != nil {
		log.Fatalf("send failed: %s", err)
	}
}

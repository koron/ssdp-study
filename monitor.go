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

func main() {
	addr, err := net.ResolveUDPAddr("udp4", addrIP4)
	if err != nil {
		log.Fatalf("net.ResolveUDPAddr() failed: %s", err)
	}
	l, err := net.ListenMulticastUDP("udp4", nil, addr)
	buf := make([]byte, 1024*1024)
	l.SetReadBuffer(len(buf))
	for {
		n, caddr, err := l.ReadFromUDP(buf)
		if err != nil {
			log.Fatalf("ReadFromUDP() failed: %s", err)
		}
		s := string(buf[:n])
		fmt.Printf("received: %q from %s\n", s, caddr.String())
	}
}

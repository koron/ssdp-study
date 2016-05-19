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

func response(addr *net.UDPAddr) error {
	c, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return err
	}
	n, err := c.Write([]byte("foobar\r\n\r\n"))
	if err != nil {
		return err
	}
	fmt.Printf("send %d bytes\n", n)
	return nil
}

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
		go func(addr *net.UDPAddr) {
			err := response(addr)
			if err != nil {
				log.Printf("response() failed: %s", err)
			}
		}(caddr)
	}
}

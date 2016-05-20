package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/koron/ssdp-study/udp"
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

func doResponse(addr *net.UDPAddr) {
	err := response(addr)
	if err != nil {
		log.Printf("response() failed: %s", err)
	}
}

func monitor(ifquery string) error {
	addr, err := net.ResolveUDPAddr("udp", addrIP4)
	if err != nil {
		return err
	}
	ifi, err := udp.Interface(ifquery)
	if err != nil {
		return err
	}
	l, err := net.ListenMulticastUDP("udp", ifi, addr)
	if err != nil {
		return err
	}
	log.Printf("listening multicast UDP for %s on %q (%s)", addr.String(), ifi.Name, ifquery)
	buf := make([]byte, 1024*1024)
	l.SetReadBuffer(len(buf))
	for {
		n, caddr, err := l.ReadFromUDP(buf)
		if err != nil {
			return err
		}
		s := string(buf[:n])
		log.Printf("received from %s %q", caddr.String(), s)
		//go doResponse(caddr)
	}
}

func main() {
	flag.Parse()
	ifq := "127.0.0.1"
	if flag.NArg() > 0 {
		ifq = flag.Arg(0)
	}
	err := monitor(ifq)
	if err != nil {
		log.Fatalf("monitor failed: %s", err)
	}
}

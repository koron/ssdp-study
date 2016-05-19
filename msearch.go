package main

import (
	"fmt"
	"log"
	"net"
	"time"
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
	addr, err := net.ResolveUDPAddr("udp4", addrIP4)
	if err != nil {
		log.Fatalf("net.ResolveUDPAddr() failed: %s", err)
	}

	connReq, err := net.DialUDP("udp4", nil, addr)
	if err != nil {
		log.Fatalf("net.DialUDP() failed: %s", err)
	}
	sec := 10
	defer connReq.Close()
	c := connReq

	// listen reponse unicast
	/*
		connResp, err := net.ListenUDP("udp4", connReq.LocalAddr().(*net.UDPAddr))
		if err != nil {
			log.Fatalf("net.ListenUDP() failed: %s", err)
		}
		defer connResp.Close()
		c = connResp
	*/

	// send
	msg := buildMSearch(stRootDevices, sec)
	if _, err := connReq.Write([]byte(msg)); err != nil {
		log.Fatalf("send failed: %s", err)
	}

	// wait response.
	buf := make([]byte, 1024*1024)
	c.SetReadBuffer(len(buf))
	c.SetReadDeadline(time.Now().Add(time.Second * time.Duration(sec)))
	for {
		n, addr, err := c.ReadFrom(buf)
		if err != nil {
			if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
				break
			}
			log.Fatalf("ReadFrom() failed: %s", err)
		}
		fmt.Printf("received from:%s %q\n", addr.String(), string(buf[:n]))
	}
}

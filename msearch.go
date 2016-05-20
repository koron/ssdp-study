package main

import (
	"flag"
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

func buildMSearch(st string, mx uint) string {
	return fmt.Sprintf("M-SEARCH * HTTP/1.1\r\nHOST: %s\r\nMAN: \"ssdp:discover\"\r\nMX: %d\r\nST: %s\r\n\r\n", addrIP4, mx, st)
}

func msearch(localAddr string, sec uint) error {
	laddr, err := net.ResolveUDPAddr("udp", localAddr)
	if err != nil {
		return err
	}
	c, err := net.ListenUDP("udp", laddr)
	if err != nil {
		return err
	}
	defer c.Close()

	fmt.Printf("local addr: %s\n", c.LocalAddr().String())

	// send
	msg := buildMSearch(stAll, sec)
	raddr, err := net.ResolveUDPAddr("udp", addrIP4)
	if err != nil {
		return err
	}
	if _, err := c.WriteTo([]byte(msg), raddr); err != nil {
		return err
	}

	// wait response.
	buf := make([]byte, 65535)
	c.SetReadBuffer(len(buf))
	c.SetReadDeadline(time.Now().Add(time.Second * time.Duration(sec)))
	for {
		n, addr, err := c.ReadFrom(buf)
		if err != nil {
			if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
				break
			}
			return err
		}
		fmt.Printf("received from:%s %q\n", addr.String(), string(buf[:n]))
	}

	return nil
}

func main() {
	sec := flag.Uint("mx", 1, "seconds to wait response")
	flag.Parse()
	if *sec < 1 {
		*sec = 1
	}
	var localAddr string
	if flag.NArg() > 0 {
		localAddr = flag.Arg(0)
	}
	err := msearch(localAddr, *sec)
	if err != nil {
		log.Fatalf("cast failed: %s", err)
	}
}

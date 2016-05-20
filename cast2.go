package main

import (
	"flag"
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

func cast(localAddr string) error {
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
	msg := buildMSearch(stAll, 1)
	raddr, err := net.ResolveUDPAddr("udp", addrIP4)
	if err != nil {
		return err
	}
	if _, err := c.WriteTo([]byte(msg), raddr); err != nil {
		return err
	}
	return nil
}

func main() {
	flag.Parse()
	var localAddr string
	if flag.NArg() > 0 {
		localAddr = flag.Arg(0)
	}
	err := cast(localAddr)
	if err != nil {
		log.Fatalf("cast failed: %s", err)
	}
}

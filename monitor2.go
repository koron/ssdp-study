package main

import (
	"errors"
	"log"
	"net"

	"github.com/koron/ssdp-study/udp"
	"golang.org/x/net/ipv4"
)

func monitor() error {
	group, err := net.ResolveUDPAddr("udp4", "239.255.255.250:1900")
	if err != nil {
		return nil
	}
	iflist, err := net.Interfaces()
	if err != nil {
		return err
	}

	c, err := net.ListenPacket("udp4", "0.0.0.0:1900")
	if err != nil {
		return err
	}
	defer c.Close()
	log.Printf("listening %s", c.LocalAddr().String())

	p := ipv4.NewPacketConn(c)
	defer p.Close()
	p.SetMulticastLoopback(true)
	empty := true
	for _, ifi := range iflist {
		if !udp.HasRealAddress(&ifi) {
			continue
		}
		p.JoinGroup(&ifi, group)
		empty = false
		log.Printf("%q joined group %s", ifi.Name, group.String())
	}
	if empty {
		return errors.New("no interfaces to listen")
	}

	buf := make([]byte, 65535)
	for {
		n, addr, err := c.ReadFrom(buf)
		if err != nil {
			return err
		}
		s := string(buf[:n])
		log.Printf("received from %s %q", addr.String(), s)
	}

	return nil
}

func main() {
	if err := monitor(); err != nil {
		log.Fatalf("monitor failed: %s", err)
	}
}

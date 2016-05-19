package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	ini, err := net.Interfaces()
	if err != nil {
		log.Fatalf("failed to obtain interfaces: %s", err)
	}
	for i, ifi := range ini {
		fmt.Printf("%d: %s\n", i, ifi.Name)
		if addrs, err := ifi.Addrs(); err == nil {
			fmt.Printf("  index: %d\n", ifi.Index)
			fmt.Println("  addrs:")
			for _, addr := range addrs {
				fmt.Printf("    %s\n", addr.String())
			}
		}
		if addrs, err := ifi.MulticastAddrs(); err == nil {
			fmt.Println("  multicast-addrs:")
			for _, addr := range addrs {
				fmt.Printf("    %s\n", addr.String())
			}
		}
	}
}

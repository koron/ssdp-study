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
	for i, n := range ini {
		fmt.Printf("%d: %#v\n", i, n)
	}
}

package main

import (
	"fmt"
	"net"
)

func main() {
	addrs, err := net.LookupTXT("github.com")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(addrs)
}

package go_net

import (
	"fmt"
	"net"
)

func Run() {
	addrs, err := net.LookupTXT("github.com")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(addrs)
}

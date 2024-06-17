package netStudy

import (
	"fmt"
	"net"
)

func Init() {
	addrs, err := net.LookupHost("localhost")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(addrs)

}

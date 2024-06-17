package netStudy

import (
	"fmt"
	"net"
)

func Init() {
	addrs, err := net.LookupTXT("github.com")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(addrs)

}

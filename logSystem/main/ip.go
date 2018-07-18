package main

import (
	"net"
	"fmt"
)

var (
	localIPArray []string
)

func init() {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(fmt.Sprintf("get addrs failed, err:%v", err))

	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				localIPArray = append(localIPArray,ipnet.IP.String())
			}
		}
	}
	fmt.Println(localIPArray)
}

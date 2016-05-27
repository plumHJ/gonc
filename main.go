// nc project main.go
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/plumhj/nc/client"
	"github.com/plumhj/nc/server"
)

func main() {

	var isServerMode bool
	var isUdp bool

	flag.BoolVar(&isServerMode, "l", false, "server or not")
	flag.BoolVar(&isUdp, "u", false, "whether udp or tcp")

	flag.Parse()

	addr := "0.0.0.0"
	port := "0"

	if isServerMode {

		if len(flag.Args()) < 1 {
			fmt.Println("Usage : nc [-u] -l [binding_address] listeing_port")
			os.Exit(1)
		}

		if len(flag.Args()) >= 2 {
			addr = flag.Arg(0)
			port = flag.Arg(1)
		} else {
			port = flag.Arg(0)
		}

		if isUdp {
			server.ListenAndServe("udp", addr, port)
		} else {
			server.ListenAndServe("tcp", addr, port)
		}

	} else {

		if len(flag.Args()) < 2 {
			fmt.Println("Usage : nc [-u] address port")
			os.Exit(1)
		}

		addr = flag.Arg(0)
		port = flag.Arg(1)

		if isUdp {
			client.ConnectAndWork("udp", addr, port)
		} else {
			client.ConnectAndWork("tcp", addr, port)
		}

	}

}

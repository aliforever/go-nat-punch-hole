package main

import (
	"fmt"
	"time"

	"github.com/aliforever/go-nat-punch-hole/libs/client"

	"github.com/aliforever/go-nat-punch-hole/args"
	"github.com/aliforever/go-nat-punch-hole/libs/server"
)

func main() {
	if args.Type == "" {
		panic("please specify type by --type=server/client")
	}

	if args.Type == "client" && args.SignalingServerAddress == "" {
		panic("please specify signaling server address using --server=ip/domain")
	}

	if args.Type == "server" {
		udpServer := server.NewUDP(args.LocalAddress)
		fmt.Println("Starting listening on", args.LocalAddress)
		err := udpServer.Listen()
		if err != nil {
			panic(err)
		}
	}

	if args.Type == "client" {
		udpClient := client.NewUDP(args.SignalingServerAddress, args.LocalAddress)
		err := udpClient.Start()
		if err != nil {
			panic(fmt.Sprintf("Error in udp client: %s", err))
		}
		fmt.Println("here now")
		ok, err := udpClient.Register("test_room")
		fmt.Println(ok, err)
	}

	for true {
		time.Sleep(time.Second * 1)
	}
}
